CREATE TABLE IF NOT EXISTS equipment (
    id TEXT PRIMARY KEY,
    name VARCHAR(125) NOT NULL,
    category VARCHAR(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS lines (
    id TEXT PRIMARY KEY,
    name VARCHAR(125) NOT NULL
);

CREATE TABLE IF NOT EXISTS specifications (
    id TEXT PRIMARY KEY,
    code VARCHAR(25) NOT NULL,
    description VARCHAR(255) NOT NULL,
    deprecated BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS workshifts (
    id TEXT PRIMARY KEY,
    name VARCHAR(125) NOT NULL,
    status VARCHAR(24) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ NOT NULL,
    duration BIGINT NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_active_workshift_idx 
ON workshifts ((1)) WHERE status = 'active';

CREATE TABLE IF NOT EXISTS intervals (
    id TEXT PRIMARY KEY,
    equipment_id TEXT NOT NULL REFERENCES equipment(id),
    status VARCHAR(24) NOT NULL,
    category VARCHAR(24) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ NOT NULL,
    duration BIGINT NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_active_interval_idx 
ON intervals (equipment_id) WHERE status = 'active';

CREATE TABLE IF NOT EXISTS lines_equipment (
    line_id TEXT NOT NULL REFERENCES lines(id),
    equipment_id TEXT NOT NULL REFERENCES equipment(id),
    position INTEGER NOT NULL,
    CONSTRAINT unique_line_equipment_idx UNIQUE (line_id, equipment_id)
);

CREATE TABLE IF NOT EXISTS qualification_rules (
    equipment_id TEXT NOT NULL REFERENCES equipment(id),
    specification_id TEXT NOT NULL REFERENCES specifications(id),
    CONSTRAINT unique_qualification_rule_idx UNIQUE (equipment_id, specification_id)
);

CREATE TABLE IF NOT EXISTS downtimes_qualifications (
    downtime_id TEXT PRIMARY KEY REFERENCES intervals(id),
    specification_id TEXT NOT NULL REFERENCES specifications(id),
    comment VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS workshifts_intervals (
    workshift_id TEXT NOT NULL REFERENCES workshifts(id),
    equipment_id TEXT NOT NULL REFERENCES equipment(id),
    interval_id TEXT NOT NULL REFERENCES intervals(id),
    CONSTRAINT unique_workshift_interval_idx UNIQUE (workshift_id, interval_id)
);

CREATE INDEX IF NOT EXISTS workshift_equipment_intervals_idx 
ON workshifts_intervals (workshift_id, equipment_id);

CREATE TABLE IF NOT EXISTS production (
    interval_id TEXT NOT NULL REFERENCES intervals(id),
    category TEXT NOT NULL,
    amount BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (interval_id, category)
);

CREATE TABLE IF NOT EXISTS workshifts_production (
    workshift_id TEXT NOT NULL REFERENCES workshifts(id),
    equipment_id TEXT NOT NULL REFERENCES equipment(id),
    good BIGINT NOT NULL DEFAULT 0,
    defect BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (workshift_id, equipment_id)
);

CREATE TABLE IF NOT EXISTS oee_metrics (
    workshift_id TEXT NOT NULL REFERENCES workshifts(id),
    equipment_id TEXT NOT NULL REFERENCES equipment(id),
    availability DOUBLE PRECISION NOT NULL DEFAULT 0,
    performance DOUBLE PRECISION NOT NULL DEFAULT 0,
    quality DOUBLE PRECISION NOT NULL DEFAULT 0,
    oee DOUBLE PRECISION NOT NULL DEFAULT 0,
    PRIMARY KEY (workshift_id, equipment_id)
);