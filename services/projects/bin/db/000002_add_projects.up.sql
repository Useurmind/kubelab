

CREATE TABLE projects (
    id integer primary key generated always as identity,
    group_id integer,
    slug text,
    data jsonb,
    UNIQUE (group_id, slug),
    CONSTRAINT fk_group
        FOREIGN KEY(group_id)
            REFERENCES groups(id)
);