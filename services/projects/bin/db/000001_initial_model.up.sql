CREATE TABLE groups (
    id integer primary key generated always as identity,
    name text,
    data jsonb
);

CREATE TABLE projects (
    id integer primary key generated always as identity,
    group_id integer,
    assigned_group_name text,
    name text,
    slug text UNIQUE,
    CONSTRAINT fk_group
        FOREIGN KEY(group_id)
            REFERENCES groups(id)
);