CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE test (
    id uuid PRIMARY KEY NOT NULL,
    content TEXT NOT NULL
)
