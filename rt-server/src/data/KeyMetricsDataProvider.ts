import type { IDataProvider } from "../core/interfaces/dataProvider/IDataProvider";
import type { IKeyMetrics } from "../core/interfaces/keyMetrics/IKeyMetrics";

/**
 * Responsible for generating a list of dummy key metrics data.
 */
export class KeyMetricsDataProvider implements IDataProvider<IKeyMetrics> {
  private readonly data: IKeyMetrics[][] = [];
  private readonly default: number = 25;

  constructor() {
    this.generateDummyKeyMetrics(this.default);
  }

  /**
   * Retrieves an array of key metrics data.
   *
   * @returns {IKeyMetrics[]} An array of `IKeyMetrics` objects, each representing a dummy activity entry.
   */
  public getData(): IKeyMetrics[] {
    return this.data[Math.floor(Math.random() * this.data.length)];
  }

  /**
   * Generates a specified number of dummy key metrics sets.
   *
   * @param {number} count - The number of dummy key metrics sets to generate.
   */
  private generateDummyKeyMetrics(count: number): void {
    for (let i = 0; i < count; i++) {
      const item: IKeyMetrics[] = [
        {
          title: "Total Jobs",
          value: this.getRandomNumberValue(1000, 2000),
        },
        {
          title: "Active Workers",
          value: this.getRandomNumberValue(50, 100),
        },
        {
          title: "Error Rate",
          value: this.getRandomPercentageValue(),
        },
      ];
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

  /**
   * Generates a random percentage value as a string.
   *
   * @returns {string} A random percentage value.
   */
  private getRandomPercentageValue(): string {
    const percentage = (Math.random() * 5).toFixed(2); // Generates a random percentage between 0% and 5%
    return `${percentage}%`;
  }
}
