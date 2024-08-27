import type { IDataProvider } from "../core/interfaces/dataProvider/IDataProvider";
import type { IRecentActivity } from "../core/interfaces/recentActivity/IRecentActivity";

/**
 * Responsible for generating a list of dummy recent activity data.
 */
export class RecentActivityDataProvider
  implements IDataProvider<IRecentActivity>
{
  /**
   * Retrieves an array of recent activity data.
   *
   * @returns {IRecentActivity[]} An array of `RecentActivity` objects, each representing a dummy activity entry.
   */
  public getData(): IRecentActivity[] {
    return this.generateDummyActivities(25);
  }

  /**
   * Generates a specified number of dummy recent activity entries.
   *
   * @param {number} count - The number of dummy activities to generate.
   * @returns {IRecentActivity[]} An array of `RecentActivity` objects, each containing random job details.
   */
  private generateDummyActivities(count: number): IRecentActivity[] {
    const activities: IRecentActivity[] = [];

    for (let i = 0; i < count; i++) {
      const startTime: string = this.generateRandomTime();
      const endTime = Math.random() > 0.5 ? this.generateRandomTime() : "-"; // 50% chance of being in progress

      activities.push({
        jobId: (i + 1).toString(),
        status:
          endTime === "-"
            ? "In Progress"
            : Math.random() > 0.5
              ? "Completed"
              : "Failed",
        startTime: startTime,
        endTime: endTime,
      });
    }

    return activities;
  }

  /**
   * Generates a random time within a range of +/- 15 minutes from the current time.
   * Simulates realistic start and end times for activities.
   *
   * @returns {string} The generated time formatted as a string in "hh:mm AM/PM" format.
   */
  private generateRandomTime(): string {
    const now = new Date();
    const offset = Math.floor(Math.random() * 30) - 15; // Random offset between -15 and +15 minutes
    const adjustedTime = new Date(now.getTime() + offset * 60000);
    return adjustedTime.toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
    });
  }
}
