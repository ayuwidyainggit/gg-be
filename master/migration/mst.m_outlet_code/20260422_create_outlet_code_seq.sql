CREATE TABLE IF NOT EXISTS mst.m_outlet_code_seq (
    outlet_code_id UUID PRIMARY KEY,
    last_sequence_no INTEGER NOT NULL DEFAULT 0,
    updated_by VARCHAR(50),
    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
);
