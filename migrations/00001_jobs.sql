-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS jobs (
  id BIGSERIAL PRIMARY KEY,
  status TEXT NOT NULL,
  data JSONB NOT NULL, 
  attempts SMALLINT NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS jobs;
-- +goose StatementEnd
