/*======================================================================*/
--  db_setup.sql
--   -- :mode=pl-sql:tabSize=3:indentSize=3:
--  Mon Aug 17 14:44:44 PST 2015 @144 /Internet Time/
--  Purpose:
--  NOTE: must be connected as 'postgres' user or a superuser to start.
/*======================================================================*/

\set ON_ERROR_STOP on
set client_min_messages to 'warning';



-- Create admin user
CREATE USER iotuser WITH PASSWORD 'dev';
CREATE DATABASE iotdb;
GRANT ALL PRIVILEGES ON DATABASE iotdb to iotuser;
ALTER USER iotuser WITH SUPERUSER;


\connect iotdb


-- Enable pgcrypto for passwords
CREATE EXTENSION IF NOT EXISTS pgcrypto;


-- @function update_modified_column
-- @description updates record updated_at column
--              with current timestamp
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = now();
        RETURN NEW;
    END;
$$ language 'plpgsql';



--
-- CONFIG
--

DROP TABLE IF EXISTS config CASCADE;

CREATE TABLE IF NOT EXISTS config (
    key             VARCHAR(50) PRIMARY KEY,
    value           VARCHAR(50),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE config IS 'Configuration and system state (key:value)';

DROP TRIGGER IF EXISTS config_update ON config;
CREATE TRIGGER config_update
    BEFORE UPDATE ON config
        FOR EACH ROW
            EXECUTE PROCEDURE update_modified_column();

INSERT INTO config(key, value) VALUES('version', '0.0.1');



--
-- DEVICES
--

DROP TABLE IF EXISTS devices CASCADE;

CREATE TABLE IF NOT EXISTS devices (
    id              VARCHAR(36) PRIMARY KEY DEFAULT md5(random()::text || now()::text)::uuid,
    name            VARCHAR(50) NOT NULL CHECK(name != ''),
    apikey          VARCHAR(32) NOT NULL UNIQUE DEFAULT md5(random()::text),
    secret_token    VARCHAR(32) NOT NULL DEFAULT md5(random()::text),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted      BOOLEAN DEFAULT false,
    is_active       BOOLEAN DEFAULT true
);

COMMENT ON TABLE devices IS 'Devices for data collection and queries';
COMMENT ON COLUMN devices.name IS 'Name of the device';

DROP TRIGGER IF EXISTS devices_update ON devices;
CREATE TRIGGER device_update
    BEFORE UPDATE ON devices
        FOR EACH ROW
            EXECUTE PROCEDURE update_modified_column();





--
-- CATEGORIES
--

DROP TABLE IF EXISTS device_attributes CASCADE;

CREATE TABLE IF NOT EXISTS device_attributes (
    device_id       VARCHAR(36) NOT NULL CHECK(device_id != ''),
    key             VARCHAR(50),
    value           VARCHAR(50),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);

DROP TRIGGER IF EXISTS device_attributes_update ON device_attributes;
CREATE TRIGGER device_attributes_update
    BEFORE UPDATE ON device_attributes
        FOR EACH ROW
            EXECUTE PROCEDURE update_modified_column();




--
-- QUEUE
--

DROP TABLE IF EXISTS queue CASCADE;

CREATE TABLE IF NOT EXISTS queue (
    id              VARCHAR(36) PRIMARY KEY DEFAULT md5(random()::text || now()::text)::uuid,
    name            VARCHAR(50) NOT NULL CHECK(name != ''),
    query           JSONB,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted      BOOLEAN DEFAULT false,
    is_active       BOOLEAN DEFAULT true
);

DROP TRIGGER IF EXISTS queue_update ON queue;
CREATE TRIGGER queue_update
    BEFORE UPDATE ON queue
        FOR EACH ROW
            EXECUTE PROCEDURE update_modified_column();


--
--
--

DROP TABLE IF EXISTS job2attribute;

CREATE TABLE IF NOT EXISTS job2attribute (
    query_id    VARCHAR(36) NOT NULL CHECK(query_id != ''),
    key             VARCHAR(50),
    value           VARCHAR(50),
    FOREIGN KEY (query_id) REFERENCES queue(id) ON DELETE CASCADE
);







INSERT INTO devices (name) VALUES ('test_device');

INSERT INTO device_attributes (device_id, key, value) VALUES
    ('2e79f4a9-d085-51fe-b842-f93392e2e0a1', 'make', 'subaru'),
    ('2e79f4a9-d085-51fe-b842-f93392e2e0a1', 'model', 'forester'),
    ('2e79f4a9-d085-51fe-b842-f93392e2e0a1', 'year', '2019'),
    ('2e79f4a9-d085-51fe-b842-f93392e2e0a1', 'color', 'blue');

INSERT INTO queue (name, query) VALUES ('test_query_1', '{"method":"ping"}');
INSERT INTO queue (name, query) VALUES ('test_query_2', '{"method":"ping"}');
INSERT INTO queue (name, query) VALUES ('test_query_3', '{"method":"ping"}');

INSERT INTO job2attribute(query_id, key, value) VALUES ('d123ce0d-043c-3187-8e4f-2c97b9f2f25a', 'make', 'subaru');
INSERT INTO job2attribute(query_id, key, value) VALUES ('d123ce0d-043c-3187-8e4f-2c97b9f2f25a', 'make', 'ford');
INSERT INTO job2attribute(query_id, key, value) VALUES ('d123ce0d-043c-3187-8e4f-2c97b9f2f25a', 'model', 'forester');
INSERT INTO job2attribute(query_id, key, value) VALUES ('a411f3a9-35a3-7248-f03c-81019e088e13', 'model', '*');


WITH jobs AS (
    SELECT
        queue.id AS query_id
    FROM queue
    INNER JOIN job2attribute
        ON job2attribute.query_id = queue.id
    INNER JOIN device_attributes
        ON (device_attributes.key = job2attribute.key OR job2attribute.key = '*')
        AND (device_attributes.value = job2attribute.value OR job2attribute.value = '*')
        AND device_attributes.device_id = '2e79f4a9-d085-51fe-b842-f93392e2e0a1'
    WHERE queue.is_deleted = 'false'
    AND  queue.is_active = 'true'
    GROUP BY queue.id
)
SELECT
    json_build_object(
        'queries',
        COALESCE(
            json_agg(queue.*),
            '[]'
        )
    ) AS data
FROM queue
INNER JOIN jobs
    ON jobs.query_id = queue.id;





--
