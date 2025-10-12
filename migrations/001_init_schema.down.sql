-- Disable foreign key checks temporarily
SET session_replication_role = replica;

DROP TABLE IF EXISTS battle_formations CASCADE;
DROP TABLE IF EXISTS enhancements CASCADE;
DROP TABLE IF EXISTS unit_keywords CASCADE;
DROP TABLE IF EXISTS keywords CASCADE;
DROP TABLE IF EXISTS rules CASCADE;
DROP TABLE IF EXISTS abilities CASCADE;
DROP TABLE IF EXISTS weapons CASCADE;
DROP TABLE IF EXISTS units CASCADE;
DROP TABLE IF EXISTS factions CASCADE;
DROP TABLE IF EXISTS games CASCADE;

DROP EXTENSION IF EXISTS "pgcrypto";

-- Re-enable foreign key checks
SET session_replication_role = DEFAULT;
