/**
 * Represents the job status data returned by the API.
 *
 * @property {string} status - The status of the jobs (e.g., "Completed", "In Progress", "Failed").
 * @property {number} count - The number of jobs in this status.
 */
export interface JobStatusService {
  status: string;
  count: number;
}

/**
 * fetchJobStatusData is a mock function simulating an API call to fetch job status data.
 * It returns a promise that resolves to an array of job status data after a delay.
 *
 * @returns {Promise<JobStatusService[]>} A promise that resolves to an array of job status data.
 */
export const fetchJobStatusData = async (): Promise<JobStatusService[]> => {
  const data: JobStatusService[] = [
    { status: "Completed", count: 500 },
    { status: "In Progress", count: 150 },
    { status: "Failed", count: 5 },
  ];

  return new Promise<JobStatusService[]>((resolve) => {
    setTimeout(() => {
      resolve(data);
    }, 1000);
  });
};
