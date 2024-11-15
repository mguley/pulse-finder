-- Drop the Trigram GIN index on the `title` column, if it exists.
DROP INDEX IF EXISTS job_vacancies_title_trgm_idx;

-- Drop the Trigram GIN index on the `company` column, if it exists.
DROP INDEX IF EXISTS job_vacancies_company_trgm_idx;

-- Drop the `job_vacancies` table, if it exists.
DROP TABLE IF EXISTS job_vacancies;

-- Optionally, remove the `pg_trgm` extension, if it is no longer needed.
DROP EXTENSION IF EXISTS pg_trgm;
