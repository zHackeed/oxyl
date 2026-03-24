/* TimescaleDB schema database - Oxyl */

-- https://stackoverflow.com/questions/7624919/check-if-a-user-defined-type-already-exists-in-postgresql
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'webhook_type') THEN
        CREATE TYPE webhook_type AS ENUM (
            'WEBHOOK',
            'DISCORD',
            'SLACK'
        );
END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'notification_type') THEN
        CREATE TYPE notification_type AS ENUM (
            'COMPANY_SETTING_UPDATE',
            'COMPANY_MEMBER_UPDATE',

            'AGENT_STATUS_UPDATE',
            'AGENT_CPU_USAGE_THRESHOLD',
            'AGENT_MEMORY_USAGE_THRESHOLD',
            'AGENT_DISK_USAGE_THRESHOLD',
            'AGENT_DISK_HEALTH_THRESHOLD',
            'AGENT_NETWORK_USAGE_THRESHOLD'
        );
END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'agent_status') THEN
        CREATE TYPE agent_status AS ENUM (
            'ACTIVE',
            'ENROLLING',
            'MAINTENANCE',
            'INACTIVE'
        );
END IF;
END $$;

CREATE TABLE IF NOT EXISTS users(
    id varchar(26) PRIMARY KEY,
    email varchar(255) NOT NULL,
    password varchar(255) NOT NULL,

    name varchar(255) NOT NULL,
    surname varchar(255) NOT NULL,

    enabled boolean NOT NULL DEFAULT true,

    last_login timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Companies related information

CREATE TABLE IF NOT EXISTS companies(
    id varchar(26) PRIMARY KEY, /* SID impl */
    display_name varchar(255) NOT NULL,

    holder varchar(255) NOT NULL, /* Simple way to obtain the owner of the company and would not need to be fetched from the users <-> companies table */
    limit_nodes int NOT NULL DEFAULT 5,

    enabled boolean NOT NULL DEFAULT true,

    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_updated timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (holder) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS company_members(
    user_id varchar(26) NOT NULL,
    company_id varchar(26) NOT NULL,

    permission_bitwise int NOT NULL DEFAULT 0,  /* todo: Permission bits, will be defined better upon handling */

    PRIMARY KEY (user_id, company_id), /* A user can be a member of multiple companies and a company can have multiple users, many to many relationship. */
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (company_id) REFERENCES companies(id)
);

CREATE TABLE IF NOT EXISTS company_notification_settings(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    holder varchar(26) NOT NULL,
    webhook_type webhook_type NOT NULL DEFAULT 'WEBHOOK',
    endpoint varchar(255) NOT NULL,
    metakeys jsonb, /* This will hold the metakeys for the notification, the body that needs to be sent. */

    FOREIGN KEY (holder) REFERENCES companies(id)
);

CREATE TABLE IF NOT EXISTS company_notification_thresholds(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    holder varchar(255) NOT NULL,

    notification_type notification_type NOT NULL,
    value int NOT NULL,

    FOREIGN KEY (holder) REFERENCES companies(id)
);

--- Agents related information

CREATE TABLE IF NOT EXISTS agents(
    id varchar(26) PRIMARY KEY, /* SID impl */
    holder varchar(26) NOT NULL, /* Company ID that is the owner of the agent */

    display_name varchar(255) NOT NULL,
    registered_ip inet NOT NULL,

    status agent_status NOT NULL DEFAULT 'ENROLLING',

    system_os varchar(255),
    cpu_model varchar(255),
    total_memory bigint,
    total_disk bigint,

    last_handshake timestamp, /* This is the last time the agent has requested a JWT update and a heartbeat. */
    last_update timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (holder) REFERENCES companies(id)
);

CREATE TABLE IF NOT EXISTS agent_partition_scheme(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    agent varchar(26) NOT NULL,

    mount_point varchar(255) NOT NULL,
    total_size bigint NOT NULL,

    is_raid boolean NOT NULL DEFAULT false,
    raid_level int,

    FOREIGN KEY (agent) REFERENCES agents(id)
    UNIQUE (agent, mount_point) -- An agent can have multiple partitions, but the agent <-> mount point is unique for each agent
);

CREATE TABLE IF NOT EXISTS agent_general_metrics(
    timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    agent varchar(26) NOT NULL,
    cpu_usage float NOT NULL,
    memory_usage float NOT NULL,
);

CREATE TABLE IF NOT EXISTS agent_disk_metrics(
    timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    agent varchar(26) NOT NULL,
    mount_point varchar(255) NOT NULL,
    disk_usage float NOT NULL,
);

CREATE TABLE IF NOT EXISTS agent_physical_disk_metrics(
    timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    agent varchar(26) NOT NULL,
    disk_path varchar(32) NOT NULL,

    health_left bigint NOT NULL, /* smart info health left percentage */
    media_errors int NOT NULL, /* smart info media errors count */
);

CREATE TABLE IF NOT EXISTS agent_network_metrics(
    timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    agent varchar(26) NOT NULL,
    interface varchar(255) NOT NULL,
    rx_bytes bigint NOT NULL,
    tx_bytes bigint NOT NULL,
);

CREATE TABLE IF NOT EXISTS agent_notification_logs(
    identifier varchar(255) NOT NULL, /* This is the identifier for the notification, it can be used to correlate with the notification settings and thresholds that triggered this notification. */
    agent varchar(26) NOT NULL,
    trigger_reason notification_type NOT NULL,
    trigger_value varchar(255) NOT NULL,

    ack boolean NOT NULL DEFAULT false,
    sent_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (identifier, sent_at)
);

SELECT create_hypertable('agent_general_metrics', 'timestamp', partition_column => 'agent', chunk_time_interval => interval '1 day');
SELECT create_hypertable('agent_disk_metrics', 'timestamp', partition_column => 'agent', chunk_time_interval => interval '1 day');
SELECT create_hypertable('agent_physical_disk_metrics', 'timestamp', partition_column => 'agent', chunk_time_interval => interval '1 day');
SELECT create_hypertable('agent_network_metrics', 'timestamp', partition_column => 'agent', chunk_time_interval => interval '1 day');

