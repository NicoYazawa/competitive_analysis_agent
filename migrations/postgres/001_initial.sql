-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "vector";

-- Products table
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(500) NOT NULL,
    category VARCHAR(200),
    brand VARCHAR(200),
    description TEXT,
    image_url VARCHAR(1000),
    source_url VARCHAR(1000),
    embedding VECTOR(1536),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Competitors table
CREATE TABLE IF NOT EXISTS competitors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(500) NOT NULL,
    platform VARCHAR(100) NOT NULL,
    platform_product_id VARCHAR(200),
    current_price DECIMAL(12, 2),
    currency VARCHAR(10) DEFAULT 'USD',
    rating DECIMAL(3, 2),
    review_count INTEGER DEFAULT 0,
    seller_rating DECIMAL(3, 2),
    seller_review_count INTEGER DEFAULT 0,
    source_url VARCHAR(1000),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Price history table
CREATE TABLE IF NOT EXISTS price_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    competitor_id UUID REFERENCES competitors(id) ON DELETE CASCADE,
    price DECIMAL(12, 2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Supply chain signals table
CREATE TABLE IF NOT EXISTS supply_chain_signals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID REFERENCES products(id) ON DELETE CASCADE,
    signal_type VARCHAR(100) NOT NULL,
    source VARCHAR(100),
    severity VARCHAR(20) DEFAULT 'medium',
    title VARCHAR(500),
    description TEXT,
    raw_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Market trends table
CREATE TABLE IF NOT EXISTS market_trends (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category VARCHAR(200) NOT NULL,
    trend_keyword VARCHAR(200),
    popularity_score DECIMAL(5, 2),
    growth_rate DECIMAL(5, 2),
    recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_competitors_platform ON competitors(platform);
CREATE INDEX idx_competitors_platform_product_id ON competitors(platform, platform_product_id);
CREATE INDEX idx_price_history_competitor_id ON price_history(competitor_id);
CREATE INDEX idx_price_history_recorded_at ON price_history(recorded_at);
CREATE INDEX idx_supply_chain_signals_product_id ON supply_chain_signals(product_id);
CREATE INDEX idx_supply_chain_signals_created_at ON supply_chain_signals(created_at);
CREATE INDEX idx_market_trends_category ON market_trends(category);
CREATE INDEX idx_market_trends_recorded_at ON market_trends(recorded_at);

-- Enable TimescaleDB hypertables
SELECT create_hypertable('price_history', 'recorded_at', if_not_exists => TRUE);
SELECT create_hypertable('market_trends', 'recorded_at', if_not_exists => TRUE);
SELECT create_hypertable('supply_chain_signals', 'created_at', if_not_exists => TRUE);
