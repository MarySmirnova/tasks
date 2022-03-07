CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS labels (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    opened TIMESTAMP NOT NULL DEFAULT now(),
    closed TIMESTAMP,
    author_id INTEGER REFERENCES users(id) DEFAULT 0,
    assigned_id INTEGER REFERENCES users(id) DEFAULT 0,
    title TEXT NOT NULL,
    content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks_labels (
    task_id INTEGER NOT NULL REFERENCES tasks(id),
    label_id INTEGER NOT NULL REFERENCES labels(id)
);