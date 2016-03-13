CREATE TABLE nodes (
    id INTEGER NOT NULL PRIMARY KEY,
    user TEXT,
    timestamp TEXT,
    lon REAL NOT NULL,
    lat REAL NOT NULL,
    marked INTEGER
);

CREATE TABLE node_tags (
    ref INTEGER NOT NULL,
    key TEXT,
    value TEXT,
    marked INTEGER
);

CREATE TABLE ways (
    id INTEGER NOT NULL PRIMARY KEY,
    user TEXT,
    timestamp TEXT,
    marked INTEGER
);

CREATE TABLE way_tags (
    ref INTEGER NOT NULL,
    key TEXT,
    value TEXT,
    marked INTEGER
);

CREATE TABLE way_nodes (
    way INTEGER NOT NULL,
    num INTEGER NOT NULL,
    node INTEGER NOT NULL
);

CREATE TABLE relations (
    id INTEGER NOT NULL PRIMARY KEY,
    user TEXT,
    timestamp TEXT,
    marked INTEGER
);

CREATE TABLE relation_tags (
    ref INTEGER NOT NULL,
    key TEXT,
    value TEXT,
    marked INTEGER
);

CREATE TABLE members (
    relation INTEGER NOT NULL,
    type TEXT,
    ref INTEGER NOT NULL,
    role TEXT,
    marked INTEGER
);