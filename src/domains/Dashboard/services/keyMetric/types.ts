/**
 * Represents a key metric entry.
 *
 * @property {string} title - The title or name of the key metric.
 * @property {number | string} value - The value of the key metric, which can be a number or a string.
 */
export interface KeyMetric {
  title: string;
  value: number | string;
}

/**
 * Enum representing socket events for key metrics statistics
 */
export enum KeyMetricSocketEvents {
  NewKeyMetrics = "newKeyMetrics",
}
