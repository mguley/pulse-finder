/**
 * Represents a recent activity entry.
 *
 * @property {string} jobId - The ID of the job.
 * @property {string} status - The status of the job (e.g., "Completed", "In Progress", "Failed").
 * @property {string} startTime - The start time of the job.
 * @property {string} endTime - The end time of the job, or "-" if still in progress.
 */
export interface RecentActivity {
  jobId: string;
  status: string;
  startTime: string;
  endTime: string;
}
