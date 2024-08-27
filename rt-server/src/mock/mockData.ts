import type { RecentActivity } from "../types";

/**
 * Utility function to generate a random date within the range of +/- 15 minutes from now.
 *
 * @returns {string} - The formatted time string in "hh:mm AM/PM" format.
 */
function generateRandomTime(): string {
  const now = new Date();
  const offset = Math.floor(Math.random() * 30) - 15; // Random offset between -15 and +15 minutes
  const adjustedTime = new Date(now.getTime() + offset * 60000);
  return adjustedTime.toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
  });
}

/**
 * Generates a list of mock activities with realistic start and end times.
 *
 * @returns {RecentActivity[]} - An array of recent activity objects.
 */
export const mockActivities: RecentActivity[] = Array.from(
  { length: 20 },
  (_, index) => {
    const startTime = generateRandomTime();
    const endTime = Math.random() > 0.5 ? generateRandomTime() : "-"; // 50% chance of being in progress

    return {
      jobId: (index + 1).toString(),
      status:
        endTime === "-"
          ? "In Progress"
          : Math.random() > 0.5
            ? "Completed"
            : "Failed",
      startTime,
      endTime,
    };
  },
);
