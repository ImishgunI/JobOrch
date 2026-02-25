CREATE TYPE IF NOT EXISTS job_status AS ENUM ('PENDING', 'PROCESSING', 'SUCCEEDED', 'FAILED');

CREATE TABLE IF NOT EXISTS jobs (
    id uuid PRIMARY KEY,
    status job_status NOT NULL,
    payload jsonb NOT NULL,
    error text,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS jobs_id_idx ON jobs (id);
CREATE INDEX IF NOT EXISTS jobs_status_idx ON jobs (status);