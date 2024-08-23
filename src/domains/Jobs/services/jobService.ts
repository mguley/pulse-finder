/**
 * Represents a job vacancy.
 *
 * @property {number} id - The unique identifier for the job.
 * @property {string} title - The title of the job position.
 * @property {string} company - The company offering the job.
 * @property {string} description - A brief description of the job role.
 * @property {string} postedAt - The date when the job was posted.
 * @property {string} location - The location where the job is based.
 */
export interface JobVacancy {
  id: number;
  title: string;
  company: string;
  description: string;
  postedAt: string;
  location: string;
}

/**
 * Represents the response from the job vacancies API.
 *
 * @property {JobVacancy[]} jobs - An array of job vacancies.
 * @property {number} itemsPerPage - The number of items to display per page.
 */
export interface JobVacancyService {
  jobs: JobVacancy[];
  itemsPerPage: number;
}

/**
 * fetchJobVacancies is a mock function simulating an API call to fetch job vacancies.
 * It returns a promise that resolves to an object containing an array of job vacancies and itemsPerPage after a delay.
 *
 * @returns {Promise<JobVacancyService>} A promise that resolves to an object with job vacancies and itemsPerPage.
 */
export const fetchJobVacancies = async (): Promise<JobVacancyService> => {
  const data: JobVacancy[] = [
    {
      id: 1,
      title: "Frontend Developer",
      company: "Tech Corp",
      description: "Develop user-facing features using React.js.",
      postedAt: "2 days ago",
      location: "Remote",
    },
    {
      id: 2,
      title: "Backend Developer",
      company: "DevWorks",
      description: "Build and maintain backend services using Node.js.",
      postedAt: "3 days ago",
      location: "New York, NY",
    },
    {
      id: 3,
      title: "Fullstack Developer",
      company: "Innovate LLC",
      description: "Work on both frontend and backend of web applications.",
      postedAt: "1 week ago",
      location: "San Francisco, CA",
    },
    {
      id: 4,
      title: "UI/UX Designer",
      company: "Creative Minds",
      description: "Design user interfaces and improve user experience.",
      postedAt: "5 days ago",
      location: "Remote",
    },
    {
      id: 5,
      title: "Data Scientist",
      company: "Data Insights",
      description: "Analyze large datasets to extract actionable insights.",
      postedAt: "1 week ago",
      location: "Boston, MA",
    },
    {
      id: 6,
      title: "DevOps Engineer",
      company: "Cloud Solutions",
      description: "Implement and manage CI/CD pipelines.",
      postedAt: "2 days ago",
      location: "Austin, TX",
    },
    {
      id: 7,
      title: "Product Manager",
      company: "Productive",
      description: "Oversee product development and strategy.",
      postedAt: "4 days ago",
      location: "Seattle, WA",
    },
    {
      id: 8,
      title: "QA Engineer",
      company: "Testify",
      description: "Ensure software quality through rigorous testing.",
      postedAt: "3 days ago",
      location: "Remote",
    },
    {
      id: 9,
      title: "System Analyst",
      company: "BizTech",
      description: "Analyze and design information systems.",
      postedAt: "1 week ago",
      location: "Chicago, IL",
    },
    {
      id: 10,
      title: "Database Administrator",
      company: "DataSecure",
      description: "Manage and secure company databases.",
      postedAt: "2 weeks ago",
      location: "Dallas, TX",
    },
    {
      id: 11,
      title: "Mobile Developer",
      company: "Appify",
      description: "Develop mobile applications for Android and iOS.",
      postedAt: "4 days ago",
      location: "Los Angeles, CA",
    },
    {
      id: 12,
      title: "Security Engineer",
      company: "SecureIT",
      description: "Implement and maintain security protocols.",
      postedAt: "1 week ago",
      location: "San Diego, CA",
    },
    {
      id: 13,
      title: "Network Engineer",
      company: "NetWorld",
      description: "Design and manage network infrastructure.",
      postedAt: "2 days ago",
      location: "Miami, FL",
    },
    {
      id: 14,
      title: "AI Engineer",
      company: "AI Innovations",
      description: "Develop AI models and integrate them into products.",
      postedAt: "5 days ago",
      location: "Boston, MA",
    },
    {
      id: 15,
      title: "Blockchain Developer",
      company: "ChainWorks",
      description: "Build and maintain blockchain-based applications.",
      postedAt: "1 week ago",
      location: "Remote",
    },
  ];
  const itemsPerPage = 8; // Mocked itemsPerPage value

  return new Promise<JobVacancyService>((resolve) => {
    setTimeout(() => {
      resolve({ jobs: data, itemsPerPage: itemsPerPage });
    }, 1000);
  });
};
