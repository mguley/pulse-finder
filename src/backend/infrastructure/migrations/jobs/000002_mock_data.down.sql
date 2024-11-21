-- Remove the inserted mock job vacancies from the `job_vacancies` table.
-- This will delete all rows matching the mock data.

DELETE FROM job_vacancies
WHERE
    title IN (
              'Frontend Developer',
              'Backend Developer',
              'Fullstack Developer',
              'UI/UX Designer',
              'Data Scientist',
              'DevOps Engineer',
              'Product Manager',
              'QA Engineer',
              'System Analyst',
              'Database Administrator',
              'Mobile Developer',
              'Security Engineer',
              'Network Engineer',
              'AI Engineer',
              'Blockchain Developer'
        )
  AND company IN (
                  'Tech Corp',
                  'DevWorks',
                  'Innovate LLC',
                  'Creative Minds',
                  'Data Insights',
                  'Cloud Solutions',
                  'Productive',
                  'Testify',
                  'BizTech',
                  'DataSecure',
                  'Appify',
                  'SecureIT',
                  'NetWorld',
                  'AI Innovations',
                  'ChainWorks'
    );
