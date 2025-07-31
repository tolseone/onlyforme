-- +goose Up
CREATE TABLE documents (
                           url TEXT NOT NULL,
                           pub_date BIGINT NOT NULL,
                           fetch_time BIGINT NOT NULL,
                           text TEXT NOT NULL,
                           PRIMARY KEY (url, fetch_time)
);

CREATE INDEX idx_documents_url ON documents(url);

-- brew install golang-migrate  # MacOS
-- go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

-- migrate create -ext sql -dir internal/migrations -seq init