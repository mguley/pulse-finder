/**
 * Represents a key metric returned by the API.
 *
 * @property {string} title - The title or name of the key metric.
 * @property {number | string} value - The value of the key metric, which can be a number or a string.
 */
export interface KeyMetricService {
  title: string;
  value: number | string;
}

/**
 * fetchKeyMetrics is a mock function simulating an API call to fetch key metrics.
 * It returns a promise that resolves to an array of key metrics after a delay.
 *
 * @returns {Promise<KeyMetricService[]>} A promise that resolves to an array of key metrics.
 */
export const fetchKeyMetrics = async (): Promise<KeyMetricService[]> => {
  const data: KeyMetricService[] = [
    { title: "Total Jobs", value: 1200 },
    { title: "Active Workers", value: 80 },
    { title: "Error Rate", value: "2%" },
  ];

  return new Promise<KeyMetricService[]>((resolve) => {
    setTimeout(() => {
      resolve(data);
    }, 1000);
  });
};
