CREATE TYPE IF NOT EXISTS job_status AS ENUM ('PENDING', 'PROCESSING', 'SUCCEEDED', 'FAILED');

CREATE TABLE IF NOT EXISTS jobs (
    id uuid PRIMARY KEY,
    status job_status NOT NULL,
    payload jsonb NOT NULL,
    error text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL,
);

CREATE INDEX IF NOT EXISTS jobs_status_idx ON jobs (status);

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tgr_update_timestamp
BEFORE UPDATE ON jobs
FOR EACH ROW
EXECUTE FUNCTION update_timestamp()