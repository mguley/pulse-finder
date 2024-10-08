import type { FC, ReactNode, ReactElement } from "react";
import { createContext, useContext, useState, useEffect } from "react";
import {
  JobStatusSocketEvents,
  type JobStatus,
} from "../services/jobStatus/types";
import type { EventHandler, WebSocket } from "../services/websocket/types";
import { SocketService } from "../services/websocket/WebSocketService";
import type { EncryptionResult } from "../services/Encryption";
import { handleDecryption } from "../utils/activityDecryption";

/**
 * Represents the structure of the JobStatusChartContext.
 *
 * @property {JobStatus[] | null} jobStatusData - An array of job status data or null if not loaded.
 * @property {boolean} loading - A flag indicating whether the job status data is being loaded.
 * @property {string | null} error - An error message if the job status fails to load, otherwise null.
 */
interface JobStatusChartContextType {
  jobStatusData: JobStatus[] | null;
  loading: boolean;
  error: string | null;
}

interface JobStatusChartProviderProps {
  children: ReactNode;
}

const JobStatusChartContext = createContext<
  JobStatusChartContextType | undefined
>(undefined);

/**
 * JobStatusChartProvider component that manages and provides the state related to job status data.
 * It fetches the job status data on mount and handles loading and error states.
 *
 * @param {JobStatusChartProviderProps} props - The properties for the JobStatusChartProvider component.
 * @returns {ReactElement} The rendered JobStatusChartProvider component.
 */
export const JobStatusChartProvider: FC<JobStatusChartProviderProps> = ({
  children,
}: JobStatusChartProviderProps): ReactElement => {
  const [jobStatusData, setJobStatusData] = useState<JobStatus[] | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const socketService: WebSocket = new SocketService();

    /**
     * Handles incoming 'newJobStatuses' events from the WebSocket connection.
     *
     * @param {JobStatus} item - The new job statuses data received from the WebSocket server.
     */
    const handleNewJobStatuses = (item: JobStatus | EncryptionResult): void => {
      const data = handleDecryption<JobStatus>(item);
      setJobStatusData((): JobStatus[] => {
        return data as JobStatus[];
      });
      setLoading(false);
    };

    /**
     * Configures an event handler to listen for 'newJobStatuses' events.
     *
     * @property {string} eventName - The name of the event to listen for.
     * @property {function} handler - The function to be called when the 'newJobStatuses' event is emitted.
     */
    const eventHandler: EventHandler = {
      eventName: JobStatusSocketEvents.NewJobStatuses,
      handler: handleNewJobStatuses,
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
    <JobStatusChartContext.Provider value={{ jobStatusData, loading, error }}>
      {children}
    </JobStatusChartContext.Provider>
  );
};

/**
 * Custom hook to access the job status chart context.
 * Throws an error if used outside of a JobStatusChartProvider
 *
 * @returns {JobStatusChartContextType} The job status chart context value.
 */
export const useJobStatusChart = (): JobStatusChartContextType => {
  const context: JobStatusChartContextType | undefined = useContext(
    JobStatusChartContext,
  );
  if (!context) {
    throw new Error(
      "useJobStatusChart must be used within a JobStatusChartProvider",
    );
  }
  return context;
};
