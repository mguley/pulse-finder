-- Create the `job_vacancies` table if it does not exist.
-- The table contains fields such as:
-- - `id`: Auto-incrementing primary key (unique identifier for each job vacancy).
-- - `title`: Title of the job vacancy, e.g., "Software Engineer".
-- - `company`: The name of the company offering the job.
-- - `description`: Description of the job role and requirements.
-- - `posted_at`: Timestamp for when the job was posted. Defaults to current time if not provided.
-- - `location`: Location of the job vacancy.
-- - `version`: Version field for optimistic concurrency control. Defaults to 1.

CREATE TABLE IF NOT EXISTS job_vacancies (
    id BIGSERIAL PRIMARY KEY,       -- Unique identifier for the job vacancy.
    title TEXT NOT NULL,            -- Title of the job vacancy.
    company TEXT NOT NULL,          -- Name of the company offering the job.
    description TEXT NOT NULL,      -- Description of the job vacancy.
    posted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),  -- Timestamp when the job was posted.
    location TEXT NOT NULL,         -- Location of the job position.
    version INTEGER NOT NULL DEFAULT 1  -- Version for optimistic concurrency control.
);

-- Enable the `pg_trgm` extension to allow the use of Trigram indexes.
-- Trigram indexing is useful for efficient substring search, improving the performance of `ILIKE` and `LIKE` queries.
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Create a GIN index on the `title` column using Trigram operators.
-- This index will allow efficient substring search on the `title` column, e.g., partial matching like `ILIKE '%engineer%'`.
CREATE INDEX job_vacancies_title_trgm_idx ON job_vacancies USING GIN (title gin_trgm_ops);

-- Create a GIN index on the `company` column using Trigram operators.
-- This index will allow efficient substring search on the `company` column, e.g., partial matching like `ILIKE '%Tech%'`.
CREATE INDEX job_vacancies_company_trgm_idx ON job_vacancies USING GIN (company gin_trgm_ops);
