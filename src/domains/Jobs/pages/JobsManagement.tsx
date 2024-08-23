import { FC, ReactElement, memo } from "react";
import JobsProviders from "../components/JobsProviders";
import JobsTable from "../components/grid/JobsTable";
import SearchJobs from "../components/search/SearchJobs";

/**
 * JobsManagement component serves as the main entry point for the Jobs Management tab.
 *
 * @returns {ReactElement} The rendered JobsManagement component.
 */
const JobsManagement: FC = (): ReactElement => {
  return (
    <JobsProviders>
      <SearchJobs />
      <JobsTable />
    </JobsProviders>
  );
};

export default memo(JobsManagement);
