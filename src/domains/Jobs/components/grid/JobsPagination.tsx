import { FC, ReactElement, memo, ChangeEvent } from "react";
import { Pagination, Box } from "@mui/material";

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (event: ChangeEvent<unknown>, page: number) => void;
}

/**
 * JobsPagination component renders a pagination control to navigate between pages of job vacancies.
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
        src={{ mt: 2 }}
      />
    </Box>
  );
};

export default memo(JobsPagination);
