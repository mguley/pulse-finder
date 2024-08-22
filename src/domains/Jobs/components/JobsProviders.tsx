import { FC, ReactElement, ReactNode, memo } from "react";
import { JobsTableProvider } from "../context/JobsTableContext";

interface JobsProvidersProps {
  children: ReactNode;
}

/**
 * JobsProviders component consolidates multiple context providers for the Jobs.
 * It serves as a wrapper that provides all necessary contexts to its children components, ensuring
 * that they have access to the appropriate state and data.
 *
 * Contexts provided:
 * - JobsTableContext: Manages the state and data for jobs grid.
 *
 * @param {JobsProvidersProps} props - The props for the JobsProviders component.
 * @returns {ReactElement} The rendered component that provides the required contexts to its children.
 */
const JobsProviders: FC<JobsProvidersProps> = ({
  children,
}: JobsProvidersProps): ReactElement => {
  return <JobsTableProvider>{children}</JobsTableProvider>;
};

export default memo(JobsProviders);
