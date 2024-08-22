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
  RecentActivityService as RecentActivity,
  fetchRecentActivity,
} from "../services/recentActivityService";

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

const RecentActivityFeedContext = createContext<
  RecentActivityFeedContextType | undefined
>(undefined);

/**
 * RecentActivityFeedProvider component that manages and provides the state related to recent activities.
 * It fetches the recent activities on mount and handles loading and error states.
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
    (async () => {
      try {
        const result: RecentActivity[] = await fetchRecentActivity();
        setRecentActivities(result);
      } catch (err) {
        setError(`Failed to fetch recent activities: ${err}`);
      } finally {
        setLoading(false);
      }
    })();
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
 * Custom hook to access the recent activity feed context.
 * Throws an error if used outside of a RecentActivityFeedProvider
 *
 * @returns {RecentActivityFeedContextType} The recent activity feed context value.
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
