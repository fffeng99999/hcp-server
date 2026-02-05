-- Common Functions
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 1. Benchmarks
CREATE TABLE IF NOT EXISTS benchmarks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    algorithm VARCHAR(50) NOT NULL,
    node_count INTEGER NOT NULL,
    duration INTEGER NOT NULL,
    target_tps INTEGER,
    actual_tps DECIMAL(10,2),
    latency_p50 DECIMAL(10,4),
    latency_p90 DECIMAL(10,4),
    latency_p99 DECIMAL(10,4),
    latency_p999 DECIMAL(10,4),
    latency_avg DECIMAL(10,4),
    latency_max DECIMAL(10,4),
    latency_min DECIMAL(10,4),
    block_count INTEGER,
    transaction_count INTEGER,
    successful_tx INTEGER,
    failed_tx INTEGER,
    block_size_avg DECIMAL(10,2),
    block_propagation_time DECIMAL(10,4),
    cpu_usage_avg DECIMAL(5,2),
    cpu_usage_max DECIMAL(5,2),
    memory_usage_avg DECIMAL(10,2),
    memory_usage_max DECIMAL(10,2),
    network_in_mbps DECIMAL(10,2),
    network_out_mbps DECIMAL(10,2),
    disk_io_read DECIMAL(10,2),
    disk_io_write DECIMAL(10,2),
    view_change_count INTEGER,
    prepare_phase_latency DECIMAL(10,4),
    commit_phase_latency DECIMAL(10,4),
    status VARCHAR(20) DEFAULT 'running',
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_algorithm CHECK (algorithm IN ('tPBFT', 'Raft', 'HotStuff', 'Leios', 'HybridPBFT')),
    CONSTRAINT chk_status CHECK (status IN ('running', 'completed', 'failed', 'cancelled')),
    CONSTRAINT chk_node_count CHECK (node_count > 0 AND node_count <= 1000),
    CONSTRAINT chk_duration CHECK (duration > 0)
);

CREATE INDEX IF NOT EXISTS idx_benchmarks_algorithm ON benchmarks(algorithm);
CREATE INDEX IF NOT EXISTS idx_benchmarks_status ON benchmarks(status);
CREATE INDEX IF NOT EXISTS idx_benchmarks_created_at ON benchmarks(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_benchmarks_node_count ON benchmarks(node_count);
CREATE INDEX IF NOT EXISTS idx_benchmarks_actual_tps ON benchmarks(actual_tps DESC);
CREATE INDEX IF NOT EXISTS idx_benchmarks_algorithm_status ON benchmarks(algorithm, status);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_benchmarks_updated_at') THEN
        CREATE TRIGGER update_benchmarks_updated_at
            BEFORE UPDATE ON benchmarks
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

-- 2. Nodes
CREATE TABLE IF NOT EXISTS nodes (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255),
    address VARCHAR(255) NOT NULL,
    public_key TEXT,
    region VARCHAR(50),
    role VARCHAR(20) DEFAULT 'validator',
    status VARCHAR(20) DEFAULT 'offline',
    trust_score DECIMAL(5,2) DEFAULT 100.00,
    uptime_percentage DECIMAL(5,2),
    total_blocks_proposed INTEGER DEFAULT 0,
    total_blocks_validated INTEGER DEFAULT 0,
    cpu_usage DECIMAL(5,2),
    memory_usage DECIMAL(10,2),
    disk_usage DECIMAL(10,2),
    peers_count INTEGER DEFAULT 0,
    network_latency_avg DECIMAL(10,4),
    last_heartbeat TIMESTAMP,
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_role CHECK (role IN ('leader', 'validator', 'follower', 'observer')),
    CONSTRAINT chk_status CHECK (status IN ('online', 'offline', 'syncing', 'failed')),
    CONSTRAINT chk_trust_score CHECK (trust_score >= 0 AND trust_score <= 100)
);

CREATE INDEX IF NOT EXISTS idx_nodes_status ON nodes(status);
CREATE INDEX IF NOT EXISTS idx_nodes_role ON nodes(role);
CREATE INDEX IF NOT EXISTS idx_nodes_trust_score ON nodes(trust_score DESC);
CREATE INDEX IF NOT EXISTS idx_nodes_last_heartbeat ON nodes(last_heartbeat DESC);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_nodes_updated_at') THEN
        CREATE TRIGGER update_nodes_updated_at
            BEFORE UPDATE ON nodes
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

-- 3. Transactions
CREATE TABLE IF NOT EXISTS transactions (
    hash VARCHAR(66) NOT NULL,
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42) NOT NULL,
    amount BIGINT NOT NULL,
    gas_price BIGINT,
    gas_limit BIGINT,
    gas_used BIGINT,
    nonce BIGINT,
    block_number BIGINT,
    block_hash VARCHAR(66),
    transaction_index INTEGER,
    status VARCHAR(20) NOT NULL,
    error_message TEXT,
    submitted_at TIMESTAMP NOT NULL,
    confirmed_at TIMESTAMP,
    latency_ms DECIMAL(10,4),
    benchmark_id UUID REFERENCES benchmarks(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_status CHECK (status IN ('pending', 'confirmed', 'failed')),
    CONSTRAINT chk_amount CHECK (amount >= 0),
    CONSTRAINT chk_addresses CHECK (from_address != to_address),
    PRIMARY KEY (hash, submitted_at)
) PARTITION BY RANGE (submitted_at);

CREATE TABLE IF NOT EXISTS transactions_default PARTITION OF transactions DEFAULT;

CREATE INDEX IF NOT EXISTS idx_transactions_from ON transactions(from_address);
CREATE INDEX IF NOT EXISTS idx_transactions_to ON transactions(to_address);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_submitted_at ON transactions(submitted_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_benchmark_id ON transactions(benchmark_id);

CREATE OR REPLACE FUNCTION create_partition_if_not_exists()
RETURNS TRIGGER AS $$
DECLARE
    partition_date DATE;
    partition_name TEXT;
    start_date DATE;
    end_date DATE;
BEGIN
    partition_date := DATE_TRUNC('month', NEW.submitted_at);
    partition_name := 'transactions_' || TO_CHAR(partition_date, 'YYYY_MM');
    start_date := partition_date;
    end_date := partition_date + INTERVAL '1 month';
    
    IF NOT EXISTS (SELECT 1 FROM pg_class WHERE relname = partition_name) THEN
        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF transactions FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, end_date
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'create_transaction_partition') THEN
        CREATE TRIGGER create_transaction_partition
            BEFORE INSERT ON transactions
            FOR EACH ROW
            EXECUTE FUNCTION create_partition_if_not_exists();
    END IF;
END $$;

-- 4. Metrics
CREATE TABLE IF NOT EXISTS metrics (
    id BIGSERIAL,
    timestamp TIMESTAMP NOT NULL,
    node_id VARCHAR(50) REFERENCES nodes(id) ON DELETE CASCADE,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL(20,6) NOT NULL,
    metric_unit VARCHAR(20),
    labels JSONB,
    benchmark_id UUID REFERENCES benchmarks(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (timestamp, node_id, metric_name)
) PARTITION BY RANGE (timestamp);

CREATE TABLE IF NOT EXISTS metrics_default PARTITION OF metrics DEFAULT;

CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON metrics(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_metrics_node_id ON metrics(node_id);
CREATE INDEX IF NOT EXISTS idx_metrics_metric_name ON metrics(metric_name);
CREATE INDEX IF NOT EXISTS idx_metrics_labels ON metrics USING GIN(labels);

-- 5. Anomalies
CREATE TABLE IF NOT EXISTS anomalies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    anomaly_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    confidence_score DECIMAL(5,2) NOT NULL,
    transaction_hash VARCHAR(66),
    node_id VARCHAR(50) REFERENCES nodes(id),
    benchmark_id UUID REFERENCES benchmarks(id),
    description TEXT,
    evidence JSONB,
    status VARCHAR(20) DEFAULT 'new',
    assigned_to VARCHAR(100),
    resolved_at TIMESTAMP,
    resolution_notes TEXT,
    detected_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_anomaly_type CHECK (
        anomaly_type IN ('wash_trade', 'spoofing', 'sandwich', 'front_running', 'ddos', 'sybil')
    ),
    CONSTRAINT chk_severity CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    CONSTRAINT chk_status CHECK (status IN ('new', 'investigating', 'confirmed', 'resolved'))
);

CREATE INDEX IF NOT EXISTS idx_anomalies_type ON anomalies(anomaly_type);
CREATE INDEX IF NOT EXISTS idx_anomalies_severity ON anomalies(severity);
CREATE INDEX IF NOT EXISTS idx_anomalies_status ON anomalies(status);
CREATE INDEX IF NOT EXISTS idx_anomalies_detected_at ON anomalies(detected_at DESC);
