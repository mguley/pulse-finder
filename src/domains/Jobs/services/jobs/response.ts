/**
 * Represents a job vacancy.
 *
 * @property {number} id - The unique identifier for the job.
 * @property {string} title - The title of the job position.
 * @property {string} company - The company offering the job.
 * @property {string} description - A brief description of the job role.
 * @property {string} posted_at - The date when the job was posted.
 * @property {string} location - The location where the job is based.
 */
export interface JobVacancy {
  id: number;
  title: string;
  company: string;
  description: string;
  posted_at: string;
  location: string;
}

/**
 * Represents the response from the API endpoint.
 *
 * @property {JobVacancy[]} jobs - An array of job vacancies.
 */
export interface JobsResponse {
  jobs: JobVacancy[];
}
