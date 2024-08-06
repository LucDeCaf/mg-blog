-- +goose Up
ALTER TABLE blog_posts RENAME TO blogs;

-- +goose Down
ALTER TABLE blogs RENAME TO blog_posts;
