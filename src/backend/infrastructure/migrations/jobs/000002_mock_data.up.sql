-- Insert mock job vacancies into the `job_vacancies` table.
-- Each entry corresponds to a job vacancy from the mock data.

INSERT INTO job_vacancies (title, company, description, posted_at, location)
VALUES
    ('Frontend Developer', 'Tech Corp', 'Develop user-facing features using React.js.', NOW() - INTERVAL '2 days', 'Remote'),
    ('Backend Developer', 'DevWorks', 'Build and maintain backend services using Node.js.', NOW() - INTERVAL '3 days', 'New York, NY'),
    ('Fullstack Developer', 'Innovate LLC', 'Work on both frontend and backend of web applications.', NOW() - INTERVAL '7 days', 'San Francisco, CA'),
    ('UI/UX Designer', 'Creative Minds', 'Design user interfaces and improve user experience.', NOW() - INTERVAL '5 days', 'Remote'),
    ('Data Scientist', 'Data Insights', 'Analyze large datasets to extract actionable insights.', NOW() - INTERVAL '7 days', 'Boston, MA'),
    ('DevOps Engineer', 'Cloud Solutions', 'Implement and manage CI/CD pipelines.', NOW() - INTERVAL '2 days', 'Austin, TX'),
    ('Product Manager', 'Productive', 'Oversee product development and strategy.', NOW() - INTERVAL '4 days', 'Seattle, WA'),
    ('QA Engineer', 'Testify', 'Ensure software quality through rigorous testing.', NOW() - INTERVAL '3 days', 'Remote'),
    ('System Analyst', 'BizTech', 'Analyze and design information systems.', NOW() - INTERVAL '7 days', 'Chicago, IL'),
    ('Database Administrator', 'DataSecure', 'Manage and secure company databases.', NOW() - INTERVAL '14 days', 'Dallas, TX'),
    ('Mobile Developer', 'Appify', 'Develop mobile applications for Android and iOS.', NOW() - INTERVAL '4 days', 'Los Angeles, CA'),
    ('Security Engineer', 'SecureIT', 'Implement and maintain security protocols.', NOW() - INTERVAL '7 days', 'San Diego, CA'),
    ('Network Engineer', 'NetWorld', 'Design and manage network infrastructure.', NOW() - INTERVAL '2 days', 'Miami, FL'),
    ('AI Engineer', 'AI Innovations', 'Develop AI models and integrate them into products.', NOW() - INTERVAL '5 days', 'Boston, MA'),
    ('Blockchain Developer', 'ChainWorks', 'Build and maintain blockchain-based applications.', NOW() - INTERVAL '7 days', 'Remote');
