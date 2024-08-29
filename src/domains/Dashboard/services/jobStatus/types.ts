/**
 * Represents a job status entry.
 *
 * @property {string} status - The status of the jobs (e.g., "Completed", "In Progress", "Failed")
 * @property {number} count - The number of jobs in this status.
 */
export interface JobStatus {
  status: string;
  count: number;
}

/**
 * Enum representing socket events for job status statistics
 */
export enum JobStatusSocketEvents {
  NewJobStatuses = "newJobStatuses",
}
