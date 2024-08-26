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

export enum RecentActivitiesSocketEvents {
  NewActivity = "newActivity",
}

export interface EventHandler {
  eventName: string;
  handler: (...args: any[]) => void;
}

export interface RecentActivitySocket {
  connect(): void;
  on(eventName: string, handler: (...args: any[]) => void): void;
  off(eventName: string, handler: (...args: any[]) => void): void;
  disconnect(): void;
  isConnected(): boolean;
}
