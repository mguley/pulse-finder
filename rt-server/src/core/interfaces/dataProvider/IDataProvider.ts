/**
 * Defines the contract for data provider.
 *
 * @template T - The type of data that the provider supplies.
 */
export interface IDataProvider<T> {
  /**
   * Retrieves or generates data.
   *
   * This method should return data of type `T`. The source of the data could vary,
   * including in-memory data, data fetched from a database, or generated dummy data.
   *
   * @returns {T} Provided data.
   */
  getData(): T[];
}
