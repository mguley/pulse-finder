CREATE TABLE IF NOT EXISTS job_vacancies (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    company TEXT NOT NULL,
    description TEXT NOT NULL,
    posted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    location TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1
);

-- Create a GIN index for full-text search on title and company
CREATE INDEX job_vacancies_title_company_tsv_idx ON job_vacancies USING GIN (
    to_tsvector('english', title || ' ' || company)
);
