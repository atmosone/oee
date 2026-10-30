CREATE TABLE IF NOT EXISTS workshifts (
	id TEXT PRIMARY KEY,
	name VARCHAR(64) NOT NULL,
	status VARCHAR(25) NOT NULL,
	duration BIGINT NOT NULL DEFAULT 0,
	started_at TIMESTAMPTZ NOT NULL,
	finished_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS equipment (
	id TEXT PRIMARY KEY,
	name VARCHAR(64) NOT NULL,
	units VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS equipment_lines (
	id TEXT PRIMARY KEY,
	name VARCHAR(64) NOT NULL,
	CONSTRAINT unique_line_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS equipment_line_items (
	equipment_line_id TEXT NOT NULL REFERENCES equipment_lines (id),
	equipment_id TEXT NOT NULL REFERENCES equipment (id),
	CONSTRAINT unique_equipment_line_item_idx UNIQUE (equipment_line_id, equipment_id)
);

CREATE TABLE IF NOT EXISTS downtimes (
	id TEXT PRIMARY KEY,
	equipment_id TEXT NOT NULL REFERENCES equipment (id),
	status VARCHAR(25),
	duration BIGINT NOT NULL DEFAULT 0,
	started_at TIMESTAMPTZ NOT NULL,
	finished_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_active_downtime_idx 
ON downtimes (equipment_id) WHERE status = 'active';

CREATE TABLE IF NOT EXISTS specifications (
	id TEXT PRIMARY KEY,
	code VARCHAR(25) NOT NULL,
	description VARCHAR(128) NOT NULL,
	CONSTRAINT unique_specification_idx UNIQUE (description)
);

CREATE TABLE IF NOT EXISTS qualification_rules (
	equipment_id TEXT NOT NULL REFERENCES equipment (id),
	specification_id TEXT NOT NULL REFERENCES specifications (id),
	deprecated BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX IF NOT EXISTS qualification_rule_idx 
ON qualification_rules (equipment_id, specification_id) 
WHERE deprecated = false;

CREATE TABLE IF NOT EXISTS qualifications (
	downtime_id TEXT PRIMARY KEY REFERENCES downtimes (id),
	specification_id TEXT NOT NULL REFERENCES specifications (id),
	comment VARCHAR(256) NOT NULL
);