ALTER TABLE mst.auto_replenishment_product
    DROP CONSTRAINT IF EXISTS auto_replenishment_product_limit_action_check;

ALTER TABLE mst.auto_replenishment_product
    DROP CONSTRAINT IF EXISTS chk_auto_replenishment_product_limit_action;

ALTER TABLE mst.auto_replenishment_product
    ADD CONSTRAINT chk_auto_replenishment_product_limit_action
    CHECK (LOWER(limit_action) IN ('restricted', 'warning', 'unrestricted'));
