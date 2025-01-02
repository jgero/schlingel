CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS resource_groups (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    resource_group_id uuid,
    FOREIGN KEY (resource_group_id) REFERENCES resource_groups(id)
);

CREATE TABLE IF NOT EXISTS resources (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    hash bytea,
    resource_group_id uuid,
    FOREIGN KEY(resource_group_id) REFERENCES resource_groups(id)
);
