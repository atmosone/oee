CREATE TABLE equipment (
    id TEXT PRIMARY KEY,
    name VARCHAR(125) NOT NULL,
    category VARCHAR(64) NOT NULL
);

CREATE TABLE lines (
    id TEXT PRIMARY KEY,
    name VARCHAR(125) NOT NULL
);

CREATE TABLE lines_equipment (
    line_id TEXT NOT NULL REFERENCES lines(id),
    equipment_id TEXT NOT NULL REFERENCES equipment(id),
    position INTEGER NOT NULL,
    CONSTRAINT unique_line_equipment_idx UNIQUE (line_id, equipment_id)
);

CREATE TABLE specifications (
    id TEXT PRIMARY KEY,
    code VARCHAR(25) NOT NULL,
    description VARCHAR(255) NOT NULL,
    deprecated BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE workshifts (
    id TEXT PRIMARY KEY,
    name VARCHAR(125) NOT NULL,
    status VARCHAR(24) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ NOT NULL,
    duration BIGINT NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_active_workshift_idx 
ON workshifts ((1)) WHERE status = 'active';

CREATE TABLE downtimes (
    id TEXT PRIMARY KEY,
    equipment_id TEXT NOT NULL REFERENCES equipment(id),
    status VARCHAR(24) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ NOT NULL,
    duration BIGINT NOT NULL DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_active_downtime_idx 
ON downtimes (equipment_id) WHERE status = 'active';

CREATE TABLE qualification_rules (
    equipment_id TEXT NOT NULL REFERENCES equipment(id),
    specification_id TEXT NOT NULL REFERENCES specifications(id),
    CONSTRAINT unique_qualification_rule_idx UNIQUE (equipment_id, specification_id)
);

CREATE TABLE downtimes_qualifications (
    downtime_id TEXT PRIMARY KEY REFERENCES downtimes(id),
    specification_id TEXT NOT NULL REFERENCES specifications(id),
    comment VARCHAR(255) NOT NULL
);

CREATE TABLE workshifts_downtimes (
    workshift_id TEXT NOT NULL REFERENCES workshifts(id),
    equipment_id TEXT NOT NULL REFERENCES equipment(id),
    downtime_id TEXT NOT NULL REFERENCES downtimes(id),
    CONSTRAINT unique_workshift_downtime_idx UNIQUE (workshift_id, downtime_id)
);

CREATE INDEX IF NOT EXISTS workshift_equipment_downtimes_idx 
ON workshifts_downtimes (workshift_id, equipment_id);

