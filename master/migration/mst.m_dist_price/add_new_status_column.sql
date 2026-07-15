ALTER TABLE "mst"."m_dist_price" 
  ADD COLUMN "status" int2 NOT NULL DEFAULT 1;

COMMENT ON COLUMN "mst"."m_dist_price"."status" IS '1: "Scheduled", 5: "Cancelled", 7: "Inactive", 10: "Published"';