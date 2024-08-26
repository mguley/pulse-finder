import {
  createContext,
  useContext,
  useState,
  FC,
  ReactNode,
  ReactElement,
  useEffect,
} from "react";
import {
  RecentActivity,
  RecentActivitiesSocketEvents,
} from "../services/recentActivity/types";
import { SocketService } from "../services/recentActivity/SocketService";

/**
 * Represents the structure of the RecentActivityFeedContextType.
 *
 * @property {RecentActivity[] | null} recentActivities - An array of recent activities or null if not loaded.
 * @property {boolean} loading - A flag indicating whether the recent activities are being loaded.
 * @property {string | null} error - An error message if the recent activities fail to load, otherwise null.
 * @property {boolean} isConnected - A flag indicating whether the socket connection is established.
 */
interface RecentActivityFeedContextType {
  recentActivities: RecentActivity[] | null;
  loading: boolean;
  error: string | null;
  isConnected: boolean;
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
  const [isConnected, setIsConnected] = useState<boolean>(false);

  useEffect(() => {
    try {
      const socketService = new SocketService(
        "https://github-io-pulse-finder.loca.lt",
      );
      socketService.connect();
      setIsConnected(socketService.isConnected());

      socketService.on(
        RecentActivitiesSocketEvents.NewActivity,
        (activity: RecentActivity): void => {
          setRecentActivities((prevState: RecentActivity[] | null) => {
            if (prevState) {
              const newState = [activity, ...prevState];
              return newState.slice(0, 3);
            }
            return [activity];
          });
        },
      );

      return () => {
        socketService.disconnect();
      };
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }, []);

  return (
    <RecentActivityFeedContext.Provider
      value={{ recentActivities, loading, error, isConnected }}
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
