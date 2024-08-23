import {
  createContext,
  useContext,
  useState,
  FC,
  ReactNode,
  ReactElement,
} from "react";
import { JobVacancy } from "../services/jobService";

/**
 * Represents the structure of the JobDetailsModalContext.
 *
 * @property {boolean} open - A flag indicating whether the modal is currently open.
 * @property {JobVacancy | null} selectedJob - The currently selected job vacancy, or null if none is selected.
 * @property {function} handleOpen - Function to open the modal with a specific job vacancy.
 * @property {function} handleClose - Function to close the modal and clear the selected job.
 */
interface JobDetailsModalContextType {
  open: boolean;
  selectedJob: JobVacancy | null;
  handleOpen: (job: JobVacancy) => void;
  handleClose: () => void;
}

interface JobDetailsModalProviderProps {
  children: ReactNode;
}

const JobDetailsModalContext = createContext<
  JobDetailsModalContextType | undefined
>(undefined);

/**
 * JobsDetailsModalProvider component manages the state for the job details modal.
 * It provides context for opening and closing the modal and for storing the currently selected job vacancy.
 *
 * @param {JobDetailsModalProviderProps} props - The properties for the JobsDetailsModalProvider component.
 * @returns {ReactElement} The rendered JobsDetailsModalProvider component.
 */
export const JobsDetailsModalProvider: FC<JobDetailsModalProviderProps> = ({
  children,
}: JobDetailsModalProviderProps): ReactElement => {
  const [open, setOpen] = useState<boolean>(false);
  const [selectedJob, setSelectedJob] = useState<JobVacancy | null>(null);

  const handleOpen = (job: JobVacancy): void => {
    setSelectedJob(job);
    setOpen(true);
  };

  const handleClose = (): void => {
    setSelectedJob(null);
    setOpen(false);
  };

  return (
    <JobDetailsModalContext.Provider
      value={{ open, selectedJob, handleOpen, handleClose }}
    >
      {children}
    </JobDetailsModalContext.Provider>
  );
};

/**
 * Custom hook to access the JobDetailsModalContext.
 * Throws an error if used outside of a JobsDetailsModalProvider.
 *
 * @returns {JobDetailsModalContextType} The context value containing modal state and functions.
 */
export const useJobDetailsModal = (): JobDetailsModalContextType => {
  const context: JobDetailsModalContextType | undefined = useContext(
    JobDetailsModalContext,
  );
  if (!context) {
    throw new Error(
      "useJobDetailsModal must be used within a JobsDetailsModalProvider",
    );
  }
  return context;
};
