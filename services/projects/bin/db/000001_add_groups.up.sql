CREATE SEQUENCE id_groups START 1;

CREATE TABLE groups (
    id integer NOT NULL DEFAULT nextval('id_groups'),
    name text NOT NULL,
    slug text NOT NULL UNIQUE,
    data jsonb
);