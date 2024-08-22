/**
 * Represents a recent activity entry returned by the API.
 *
 * @property {string} jobId - The ID of the job.
 * @property {string} status - The status of the job (e.g., "Completed", "In Progress", "Failed").
 * @property {string} startTime - The start time of the job.
 * @property {string} endTime - The end time of the job, or "-" if still in progress.
 */
export interface RecentActivityService {
  jobId: string;
  status: string;
  startTime: string;
  endTime: string;
}

/**
 * fetchRecentActivity is a mock function simulating an API call to fetch recent activity.
 * It returns a promise that resolves to an array of recent activity data after a delay.
 *
 * @returns {Promise<RecentActivityService[]>} A promise that resolves to an array of recent activity data.
 */
export const fetchRecentActivity = async (): Promise<
  RecentActivityService[]
> => {
  const data: RecentActivityService[] = [
    {
      jobId: "1234",
      status: "Completed",
      startTime: "10:00 AM",
      endTime: "10:05 AM",
    },
    {
      jobId: "1235",
      status: "In Progress",
      startTime: "10:10 AM",
      endTime: "-",
    },
    {
      jobId: "1236",
      status: "Failed",
      startTime: "10:15 AM",
      endTime: "10:20 AM",
    },
  ];

  return new Promise<RecentActivityService[]>((resolve) => {
    setTimeout(() => {
      resolve(data);
    }, 1000);
  });
};
