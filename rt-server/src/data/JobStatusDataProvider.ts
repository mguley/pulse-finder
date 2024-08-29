import type { IDataProvider } from "../core/interfaces/dataProvider/IDataProvider";
import type { IJobStatus } from "../core/interfaces/jobStatus/IJobStatus";

/**
 * Responsible for generating a list of dummy job status data.
 */
export class JobStatusDataProvider implements IDataProvider<IJobStatus> {
  private readonly data: IJobStatus[][] = [];
  private readonly defaultCount: number = 25;

  constructor() {
    this.generateDummyJobStatuses(this.defaultCount);
  }

  /**
   * Retrieves an array of jobs status data.
   *
   * @returns{IJobStatus[]} An array of `IJobStatus` objects, each representing a dummy statistic.
   */
  public getData(): IJobStatus[] {
    return this.data[Math.floor(Math.random() * this.data.length)];
  }

  /**
   * Generates a specified number of dummy job statuses.
   *
   * @param {number} count - The number of dummy job statuses to generate.
   */
  private generateDummyJobStatuses(count: number): void {
    const statuses: [string, string, string] = [
      "Completed",
      "In Progress",
      "Failed",
    ];

    for (let i = 0; i < count; i++) {
      const item: IJobStatus[] = statuses.map((status: string) => ({
        status: status,
        count: this.getRandomNumberValue(0, 500),
      }));

      this.data.push(item);
    }
  }

  /**
   * Generates a random number within the specified range.
   *
   * @param {number} min - The minimum value.
   * @param {number} max - The maximum value.
   * @returns {number} A random number within the specified range.
   */
  private getRandomNumberValue(min: number, max: number): number {
    return Math.floor(Math.random() * (max - min + 1)) + min;
  }
}
