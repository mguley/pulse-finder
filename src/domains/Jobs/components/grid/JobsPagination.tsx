import { FC, ReactElement, memo, ChangeEvent } from "react";
import { Pagination, Box } from "@mui/material";

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (event: ChangeEvent<unknown>, page: number) => void;
}

/**
 * JobsPagination component renders pagination controls for navigating through pages of job vacancies.
 * It handles the current page state and triggers the appropriate callback when the page is changed.
 *
 * @param {number} currentPage - The current page number.
 * @param {number} totalPages - The total number of pages available.
 * @param {function} onPageChange - The function to handle page changes.
 * @returns {ReactElement} The rendered component that displays pagination controls.
 */
const JobsPagination: FC<PaginationProps> = ({
  currentPage,
  totalPages,
  onPageChange,
}: PaginationProps): ReactElement => {
  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        mt: 4,
        mb: 4,
      }}
    >
      <Pagination
        page={currentPage}
        count={totalPages}
        onChange={onPageChange}
        variant="outlined"
        shape="rounded"
        color="primary"
      />
    </Box>
  );
};

export default memo(JobsPagination);
