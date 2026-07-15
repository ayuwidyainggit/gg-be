ALTER TABLE mst.m_product
ADD COLUMN origin TEXT NOT NULL DEFAULT 'create',
ADD COLUMN assigner_user_id BIGINT;

-- Constrain origin values
ALTER TABLE mst.m_product
ADD CONSTRAINT m_product_origin_chk
CHECK (origin = LOWER(origin)
   AND origin IN ('import', 'assignment', 'create', 'bulk'));

-- -- Foreign key to sys.m_user
-- -- ERROR: there is no unique constraint matching given keys for referenced table "m_user"
-- ALTER TABLE mst.m_product
-- ADD CONSTRAINT m_product_assigner_user_fk
-- FOREIGN KEY (assigner_user_id)
-- REFERENCES sys.m_user(id)
-- ON UPDATE CASCADE
-- ON DELETE SET NULL;

CREATE INDEX idx_m_product_assigner_user_id
ON mst.m_product(assigner_user_id);

-- ALTER TABLE mst.m_product
-- ALTER COLUMN origin SET DEFAULT 'create';