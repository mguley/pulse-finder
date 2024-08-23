import { FC, ReactElement, memo, ChangeEvent } from "react";
import JobCards from "./JobCards";
import JobsPagination from "./JobsPagination";
import { useJobs } from "../../context/JobsTableContext";
import LoadingSpinner from "../../../../shared/components/LoadingSpinner";
import Typography from "@mui/material/Typography";

/**
 * JobsTable component manages the display of job vacancies in a paginated grid format.
 * It retrieves job data and pagination settings from the JobsTableContext.
 * The component handles loading states, errors, and pagination controls.
 *
 * @returns {ReactElement} The rendered component that displays job vacancies in a paginated grid.
 */
const JobsTable: FC = (): ReactElement => {
  const {
    filteredJobs,
    itemsPerPage,
    currentPage,
    setCurrentPage,
    loading,
    error,
  } = useJobs();

  if (loading) {
    return <LoadingSpinner />;
  }

  if (error) {
    return (
      <Typography variant="h6" color="error">
        {error}
      </Typography>
    );
  }

  if (!filteredJobs || filteredJobs.length === 0) {
    return <Typography variant="h6">No jobs available.</Typography>;
  }

  // Calculate the total pages based on the number of jobs and items per page
  const totalPages = Math.ceil(filteredJobs.length / itemsPerPage);

  // Calculate the slice of jobs to display for the current page
  const displayedJobs = filteredJobs.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage,
  );

  const handlePageChange = (
    event: ChangeEvent<unknown>,
    page: number,
  ): void => {
    setCurrentPage(page);
  };

  return (
    <>
      <JobCards jobs={displayedJobs} />
      <JobsPagination
        currentPage={currentPage}
        totalPages={totalPages}
        onPageChange={handlePageChange}
      />
    </>
  );
};

export default memo(JobsTable);
