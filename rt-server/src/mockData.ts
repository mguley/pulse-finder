import { RecentActivity } from "./types";

export const mockActivities: RecentActivity[] = [
  {
    jobId: "1",
    status: "Completed",
    startTime: "10:00 AM",
    endTime: "10:05 AM",
  },
  {
    jobId: "2",
    status: "In Progress",
    startTime: "10:10 AM",
    endTime: "-",
  },
  {
    jobId: "3",
    status: "Failed",
    startTime: "10:15 AM",
    endTime: "10:20 AM",
  },
  {
    jobId: "4",
    status: "Completed",
    startTime: "10:30 AM",
    endTime: "10:35 AM",
  },
  {
    jobId: "5",
    status: "In Progress",
    startTime: "10:40 AM",
    endTime: "-",
  },
];
