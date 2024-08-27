import type { FC, ReactNode, ReactElement } from "react";
import { createContext, useContext, useState, useEffect } from "react";
import type { KeyMetricService as KeyMetric } from "../services/keyMetricService";
import { fetchKeyMetrics } from "../services/keyMetricService";

/**
 * Represents the structure of the KeyMetricsContext.
 *
 * @property {KeyMetric[] | null} keyMetrics - An array of key metrics or null if not loaded.
 * @property {boolean} loading - A flag indicating whether the key metrics are being loaded.
 * @property {string | null} error - An error message if the key metrics fail to load, otherwise null.
 */
interface KeyMetricsContextType {
  keyMetrics: KeyMetric[] | null;
  loading: boolean;
  error: string | null;
}

interface KeyMetricsProviderProps {
  children: ReactNode;
}

const KeyMetricsContext = createContext<KeyMetricsContextType | undefined>(
  undefined,
);

/**
 * KeyMetricsProvider component that manages and provides the state related to key metrics.
 * It fetches the key metrics on mount and handles loading and error states.
 *
 * @param {KeyMetricsProviderProps} props - The properties for the KeyMetricsProvider component.
 * @returns {ReactElement} The rendered KeyMetricsProvider component.
 */
export const KeyMetricsProvider: FC<KeyMetricsProviderProps> = ({
  children,
}: KeyMetricsProviderProps): ReactElement => {
  const [keyMetrics, setKeyMetrics] = useState<KeyMetric[] | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const result: KeyMetric[] = await fetchKeyMetrics();
        setKeyMetrics(result);
      } catch (err) {
        setError(`Failed to fetch key metrics: ${err}`);
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  return (
    <KeyMetricsContext.Provider value={{ keyMetrics, loading, error }}>
      {children}
    </KeyMetricsContext.Provider>
  );
};

/**
 * Custom hook to access the key metrics context.
 * Throws an error if used outside of a KeyMetricsProvider
 *
 * @returns {KeyMetricsContextType} The key metrics context value.
 */
export const useKeyMetrics = (): KeyMetricsContextType => {
  const context: KeyMetricsContextType | undefined =
    useContext(KeyMetricsContext);
  if (!context) {
    throw new Error("useKeyMetrics must be used within a KeyMetricsProvider");
  }
  return context;
};
