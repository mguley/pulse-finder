/**
 * Represents a job status entry.
 *
 * @property {string} status - Jobs status (e.g., "Completed", "In Progress", "Failed")
 * @property {number} count - The number of jobs in this status.
 */
export interface IJobStatus {
  status: string;
  count: number;
}
