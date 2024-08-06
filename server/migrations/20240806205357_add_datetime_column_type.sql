-- +goose Up
CREATE TABLE new_blogs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES authors(id)
);

INSERT INTO new_blogs SELECT * FROM blogs;

DROP TABLE blogs;

ALTER TABLE new_blogs RENAME TO blogs;

-- +goose Down
CREATE TABLE old_blogs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES authors(id)
);

INSERT INTO old_blogs SELECT * FROM blogs;

DROP TABLE blogs;

ALTER TABLE old_blogs RENAME TO blogs;
