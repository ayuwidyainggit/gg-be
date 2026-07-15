-- Migration: Add distributor_setup columns to mst.m_distributor
-- Description: Adds 4 boolean columns for distributor setup configuration
-- Reference: docs/Enhance Distributor _ BE.md

ALTER TABLE "mst"."m_distributor" 
  ADD COLUMN IF NOT EXISTS "allow_add_product" BOOLEAN NOT NULL DEFAULT false,
  ADD COLUMN IF NOT EXISTS "allow_edit_product" BOOLEAN NOT NULL DEFAULT false,
  ADD COLUMN IF NOT EXISTS "allow_manage_pricing" BOOLEAN NOT NULL DEFAULT false,
  ADD COLUMN IF NOT EXISTS "allow_upload_secondary_sales" BOOLEAN NOT NULL DEFAULT false;

COMMENT ON COLUMN "mst"."m_distributor"."allow_add_product" IS 'Distributor Setup: Allow distributor to add product';
COMMENT ON COLUMN "mst"."m_distributor"."allow_edit_product" IS 'Distributor Setup: Allow distributor to edit product';
COMMENT ON COLUMN "mst"."m_distributor"."allow_manage_pricing" IS 'Distributor Setup: Allow distributor to manage pricing';
COMMENT ON COLUMN "mst"."m_distributor"."allow_upload_secondary_sales" IS 'Distributor Setup: Allow distributor to upload secondary sales';
