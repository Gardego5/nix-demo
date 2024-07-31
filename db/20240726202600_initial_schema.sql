-- +goose Up
-- +goose StatementBegin
CREATE TABLE gizmos (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL
);

CREATE TABLE widgets (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  gizmo_id INTEGER NOT NULL REFERENCES gizmos (id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE widgets;

DROP TABLE gizmos;

-- +goose StatementEnd
