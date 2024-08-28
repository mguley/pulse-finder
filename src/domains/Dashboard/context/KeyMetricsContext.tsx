import type { FC, ReactNode, ReactElement } from "react";
import { createContext, useContext, useState, useEffect } from "react";
import {
  KeyMetricSocketEvents,
  type KeyMetric,
} from "../services/keyMetric/types";
import type { EventHandler, WebSocket } from "../services/websocket/types";
import { SocketService } from "../services/websocket/WebSocketService";
import type { EncryptionResult } from "../services/Encryption";
import { handleDecryption } from "../utils/activityDecryption";

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
    const socketService: WebSocket = new SocketService();

    /**
     * Handles incoming 'newKeyMetrics' events from the WebSocket connection.
     *
     * @param {KeyMetric} keyMetrics - The new key metrics data received from the WebSocket server.
     */
    const handleNewKeyMetrics = (
      keyMetrics: KeyMetric | EncryptionResult,
    ): void => {
      const data = handleDecryption<KeyMetric>(keyMetrics);
      setKeyMetrics(() => {
        return data as KeyMetric[];
      });
      setLoading(false);
    };

    /**
     * Configures an event handler to listen for 'newKeyMetrics' events.
     *
     * @property {string} eventName - The name of the event to lister for.
     * @property {function} handler - The function to be called when the 'newKeyMetrics' event is emitted.
     */
    const eventHandler: EventHandler = {
      eventName: KeyMetricSocketEvents.NewKeyMetrics,
      handler: handleNewKeyMetrics,
    };

    try {
      socketService.connect();
      socketService.on(eventHandler);

      return () => {
        socketService.off(eventHandler);
        socketService.disconnect();
      };
    } catch (e) {
      setError((e as Error)?.message || "An unknown error occurred");
    }
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
