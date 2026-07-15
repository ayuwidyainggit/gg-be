-- Optional promotion header fields: budget id, claim date range, VAT/WHT rates
ALTER TABLE promo.promotions
  ADD COLUMN IF NOT EXISTS budget_id varchar(50),
  ADD COLUMN IF NOT EXISTS claim_date_from date,
  ADD COLUMN IF NOT EXISTS claim_date_to date,
  ADD COLUMN IF NOT EXISTS vat_rate numeric(5,2),
  ADD COLUMN IF NOT EXISTS wht_rate numeric(5,2);

COMMENT ON COLUMN promo.promotions.budget_id IS 'External budget identifier (alphanumeric)';
COMMENT ON COLUMN promo.promotions.claim_date_from IS 'Claim period start date';
COMMENT ON COLUMN promo.promotions.claim_date_to IS 'Claim period end date';
COMMENT ON COLUMN promo.promotions.vat_rate IS 'VAT rate percentage (0-100)';
COMMENT ON COLUMN promo.promotions.wht_rate IS 'WHT rate percentage (0-100)';
