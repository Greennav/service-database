CREATE TABLE IF NOT EXISTS nodes (
    id INTEGER NOT NULL PRIMARY KEY,
    user TEXT,
    timestamp TEXT,
    lon REAL NOT NULL,
    lat REAL NOT NULL,
    marked INTEGER
);

CREATE TABLE IF NOT EXISTS node_tags (
    ref INTEGER NOT NULL,
    key TEXT,
    value TEXT,
    marked INTEGER
);

CREATE TABLE IF NOT EXISTS ways (
    id INTEGER NOT NULL PRIMARY KEY,
    user TEXT,
    timestamp TEXT,
    marked INTEGER
);

CREATE TABLE IF NOT EXISTS way_tags (
    ref INTEGER NOT NULL,
    key TEXT,
    value TEXT,
    marked INTEGER
);

CREATE TABLE IF NOT EXISTS way_nodes (
    way INTEGER NOT NULL,
    num INTEGER NOT NULL,
    node INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS relations (
    id INTEGER NOT NULL PRIMARY KEY,
    user TEXT,
    timestamp TEXT,
    marked INTEGER
);

CREATE TABLE IF NOT EXISTS relation_tags (
    ref INTEGER NOT NULL,
    key TEXT,
    value TEXT,
    marked INTEGER
);

CREATE TABLE IF NOT EXISTS members (
    relation INTEGER NOT NULL,
    type TEXT,
    ref INTEGER NOT NULL,
    role TEXT,
    marked INTEGER
);