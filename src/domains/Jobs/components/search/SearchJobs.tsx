import { FC, ReactElement, memo, ChangeEvent } from "react";
import { TextField } from "@mui/material";
import { useSearchJobs } from "../../context/SearchJobsContextType";

/**
 * SearchJobs component that provides an input field to search jobs by title or company.
 * The search term is managed by the SearchJobsContext.
 *
 * @returns {ReactElement} The rendered SearchJobs component.
 */
const SearchJobs: FC = (): ReactElement => {
  const { searchTerm, setSearchTerm } = useSearchJobs();
  const handleOnChange = (event: ChangeEvent<HTMLInputElement>): void => {
    setSearchTerm(event.target.value);
  };

  return (
    <TextField
      label="Search Jobs: Title, Company"
      variant="outlined"
      fullWidth
      value={searchTerm}
      onChange={handleOnChange}
      sx={{ mb: 2 }}
    />
  );
};

export default memo(SearchJobs);
