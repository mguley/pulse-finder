/**
 * Represents a recent activity item.
 *
 * @property {string} jobId - The ID of the job.
 * @property {string} status - The status of the job (e.g., "Completed", "In Progress", "Failed").
 * @property {string} startTime - The start time of the job.
 * @property {string} endTime - The end time of the job, or "-" if still in progress.
 */
export interface IRecentActivity {
  jobId: string;
  status: string;
  startTime: string;
  endTime: string;
}
