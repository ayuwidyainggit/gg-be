DB verification - tidak perlu ALTER/migration.

Query: \d mst.m_distributor + pg_constraint + pg_indexes + live INSERT (rollback)

Kolom distributor_code:
- Type: character varying(20)
- Nullable: YES
- Default: NULL
- UNIQUE: TIDAK ada di kolom ini
- NOT NULL: TIDAK
- CHECK constraint regex: TIDAK ada

Constraint di tabel mst.m_distributor (lengkap):
- fk_m_distributor_parent_cust_id (FOREIGN KEY parent_cust_id -> smc.m_customer)
- m_distributor_cust_id_not_null (NOT NULL cust_id)
- m_distributor_distributor_id_not_null (NOT NULL distributor_id)
- m_distributor_pkey (PRIMARY KEY distributor_id)

Index:
- m_distributor_pkey UNIQUE (distributor_id)
- idx_m_distributor_parent_cust_id (parent_cust_id)
- uq_m_distributor_distributor UNIQUE (distributor_id)

Live insert test (TEMP table LIKE mst.m_distributor, ROLLBACK):
INSERT distributor_code='DIST-15676761A' -> success (returned distributor_id=129).
Kesimpulan: DB sudah alphanumeric-ready. Tidak ada ALTER.

Foreign keys yang refer ke tabel ini semua via distributor_id, bukan distributor_code.
Jadi perubahan tipe/panjang distributor_code tidak akan break FK manapun.
