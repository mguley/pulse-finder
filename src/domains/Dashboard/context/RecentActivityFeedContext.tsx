import type { FC, ReactNode, ReactElement } from "react";
import { createContext, useContext, useState, useEffect } from "react";
import {
  RecentActivitiesSocketEvents,
  type RecentActivity,
} from "../services/recentActivity/types";
import type { EventHandler, WebSocket } from "../services/websocket/types";
import { SocketService } from "../services/websocket/WebSocketService";
import type { EncryptionResult } from "../services/Encryption";
import { handleDecryption } from "../utils/activityDecryption";

/**
 * Represents the structure of the RecentActivityFeedContextType.
 *
 * @property {RecentActivity[] | null} recentActivities - An array of recent activities or null if not loaded.
 * @property {boolean} loading - A flag indicating whether the recent activities are being loaded.
 * @property {string | null} error - An error message if the recent activities fail to load, otherwise null.
 */
interface RecentActivityFeedContextType {
  recentActivities: RecentActivity[] | null;
  loading: boolean;
  error: string | null;
}

interface RecentActivityFeedProviderProps {
  children: ReactNode;
}

/**
 * Context to manage the state and connection for recent activities.
 *
 * This context provider is responsible for managing the WebSocket connection,
 * fetching real-time updates for recent activities, and providing these updates
 * to the components that consume this context.
 */
const RecentActivityFeedContext = createContext<
  RecentActivityFeedContextType | undefined
>(undefined);

/**
 * RecentActivityFeedProvider component that manages and provides the state related to recent activities.
 * It connects to the WebSocket server on mount and handles loading and error states.
 *
 * @param {RecentActivityFeedProviderProps} props - The properties for the RecentActivityFeedProvider component.
 * @returns {ReactElement} The rendered RecentActivityFeedProvider component.
 */
export const RecentActivityFeedProvider: FC<
  RecentActivityFeedProviderProps
> = ({ children }: RecentActivityFeedProviderProps): ReactElement => {
  const [recentActivities, setRecentActivities] = useState<
    RecentActivity[] | null
  >(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    /**
     * Initializes the WebSocket service to establish a connection to the server.
     */
    const socketService: WebSocket = new SocketService();

    /**
     * Handles incoming 'NewActivity' events from the WebSocket connection.
     *
     * @param {RecentActivity} activity - The new activity data received from the WebSocket server.
     */
    const handleNewActivity = (
      activity: RecentActivity | EncryptionResult,
    ): void => {
      const data: RecentActivity = handleDecryption(activity);
      setRecentActivities((prevState: RecentActivity[] | null) => {
        return prevState ? [data, ...prevState].slice(0, 3) : [data];
      });
      setLoading(false);
    };

    /**
     * Configures an event handler to listen for 'newActivity' events.
     *
     * @property {string} eventName - The name of the event to listen for.
     * @property {function} handler - The function to be called when the 'newActivity' event is emitted.
     */
    const eventHandler: EventHandler = {
      eventName: RecentActivitiesSocketEvents.NewActivity,
      handler: handleNewActivity,
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
    <RecentActivityFeedContext.Provider
      value={{ recentActivities, loading, error }}
    >
      {children}
    </RecentActivityFeedContext.Provider>
  );
};

/**
 * Custom hook to access the RecentActivityFeedContext.
 *
 * This hook provides a simple way for components to access the recent activity feed data,
 * including the list of recent activities, loading state, connection status, and any errors
 * that may have occurred during the connection process.
 *
 * @returns {RecentActivityFeedContextType} The recent activity feed context value.
 * @throws Will throw an error if used outside of a RecentActivityFeedProvider.
 */
export const useRecentActivityFeed = (): RecentActivityFeedContextType => {
  const context: RecentActivityFeedContextType | undefined = useContext(
    RecentActivityFeedContext,
  );
  if (!context) {
    throw new Error(
      "useRecentActivityFeed must be used within a RecentActivityFeedProvider",
    );
  }
  return context;
};
