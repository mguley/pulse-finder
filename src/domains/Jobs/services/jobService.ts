/**
 * Represents a job vacancy.
 *
 * @property {number} id - The unique identifier for the  job.
 * @property {string} title - The title of the job position.
 * @property {string} company - The company offering the job.
 */
export interface JobVacancy {
  id: number;
  title: string;
  company: string;
}

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
    { id: 1, title: "Frontend Developer", company: "Tech Corp" },
    { id: 2, title: "Backend Developer", company: "DevWorks" },
    { id: 3, title: "Fullstack Developer", company: "Innovate LLC" },
    { id: 4, title: "UI/UX Designer", company: "Creative Minds" },
    { id: 5, title: "Data Scientist", company: "Data Insights" },
    { id: 6, title: "DevOps Engineer", company: "Cloud Solutions" },
    { id: 7, title: "Product Manager", company: "Productive" },
    { id: 8, title: "QA Engineer", company: "Testify" },
    { id: 9, title: "System Analyst", company: "BizTech" },
    { id: 10, title: "Database Administrator", company: "DataSecure" },
    { id: 11, title: "Mobile Developer", company: "Appify" },
    { id: 12, title: "Security Engineer", company: "SecureIT" },
    { id: 13, title: "Network Engineer", company: "NetWorld" },
    { id: 14, title: "AI Engineer", company: "AI Innovations" },
    { id: 15, title: "Blockchain Developer", company: "ChainWorks" },
  ];
  const itemsPerPage = 8; // Mocked itemsPerPage value

  return new Promise<JobVacancyService>((resolve) => {
    setTimeout(() => {
      resolve({ jobs: data, itemsPerPage: itemsPerPage });
    }, 1000);
  });
};
