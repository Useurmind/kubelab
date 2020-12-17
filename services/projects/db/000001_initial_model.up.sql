CREATE TABLE groups (
    id integer primary key generated always as identity,
    name text,
    is_root boolean
);

CREATE TABLE groups2groups (
    id integer primary key generated always as identity,
    parent_id integer,
    child_id integer
    CONSTRAINT fk_parent_group
      FOREIGN KEY(parent_id) 
	  REFERENCES groups(id)

    CONSTRAINT fk_child_group
      FOREIGN KEY(child_id) 
	  REFERENCES groups(id)
);

CREATE TABLE projects (
    id integer primary key generated always as identity,
    name text,
    slug text UNIQUE
);