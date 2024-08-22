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
  JobVacancy,
  JobVacancyService,
  fetchJobVacancies,
} from "../services/jobService";

/**
 * Represents the structure of the JobsTableContext.
 *
 * @property {JobVacancy[] | null} jobs - An array of job vacancies or null if not loaded.
 * @property {number} itemsPerPage - The number of items to display per page.
 * @property {number} currentPage - The current page number.
 * @property {boolean} loading - A flag indicating whether the job vacancies are being loaded.
 * @property {string | null} error - An error message if the job vacancies fail to load, otherwise null.
 * @property {function} setCurrentPage - Function to set the current page number.
 */
interface JobsTableContextType {
  jobs: JobVacancy[] | null;
  itemsPerPage: number;
  currentPage: number;
  loading: boolean;
  error: string | null;
  setCurrentPage: (page: number) => void;
}

interface JobsTableProviderProps {
  children: ReactNode;
}

const JobsTableContext = createContext<JobsTableContextType | undefined>(
  undefined,
);

/**
 * JobsTableProvider component that manages and provides the state related to job vacancies.
 * It fetches the job vacancies on mount and handles loading and error states.
 *
 * @param {JobsTableProviderProps} props - The properties for the JobsTableProvider component.
 * @returns {ReactElement} The rendered JobsTableProvider component.
 */
export const JobsTableProvider: FC<JobsTableProviderProps> = ({
  children,
}: JobsTableProviderProps): ReactElement => {
  const [jobs, setJobs] = useState<JobVacancy[] | null>(null);
  const [itemsPerPage, setItemsPerPage] = useState<number>(3);
  const [currentPage, setCurrentPage] = useState(1);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    (async () => {
      try {
        const result: JobVacancyService = await fetchJobVacancies();
        setJobs(result.jobs);
        setItemsPerPage(result.itemsPerPage);
      } catch (err) {
        setError(`Failed to load jobs data: ${err}`);
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  return (
    <JobsTableContext.Provider
      value={{
        jobs,
        itemsPerPage,
        currentPage,
        loading,
        error,
        setCurrentPage,
      }}
    >
      {children}
    </JobsTableContext.Provider>
  );
};

/**
 * Custom hook to access the jobs table context.
 * Throws an error if used outside of a JobsTableProvider
 *
 * @returns {JobsTableContextType} The context value containing jobs, loading and error state.
 */
export const useJobs = (): JobsTableContextType => {
  const context: JobsTableContextType | undefined =
    useContext(JobsTableContext);
  if (!context) {
    throw new Error("useJobs must be used within a JobsTableProvider");
  }
  return context;
};
