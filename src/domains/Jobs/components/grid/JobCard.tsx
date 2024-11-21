import type { FC, ReactElement } from "react";
import { memo } from "react";
import { Card, CardContent, Typography } from "@mui/material";
import type { JobVacancy } from "../../services/jobs/response";
import { useJobDetailsModal } from "../../context/JobDetailsModalContext";
import JobDetailsModal from "./JobDetailsModal";

interface JobCardProps {
  job: JobVacancy;
}

/**
 * JobCard component is a presentational component that displays the details of a single job vacancy.
 * On clicking the card, a modal pops us with detailed information about the job.
 *
 * @param {JobVacancy} job - The job vacancy object containing details to display.
 * @returns {ReactElement} The rendered component that displays a job vacancy.
 */
const JobCard: FC<JobCardProps> = ({ job }: JobCardProps): ReactElement => {
  const { open, selectedJob, handleOpen, handleClose } = useJobDetailsModal();

  return (
    <>
      <Card
        variant="outlined"
        sx={{
          transition: "background-color 0.3s",
          "&:hover": {
            backgroundColor: "#f5f5f5",
            cursor: "pointer",
          },
        }}
        onClick={() => handleOpen(job)}
      >
        <CardContent>
          <Typography variant="h6" component="div">
            {job.title}
          </Typography>
          <Typography color="text.secondary">{job.company}</Typography>
        </CardContent>
      </Card>

      {selectedJob && (
        <JobDetailsModal job={selectedJob} open={open} onClose={handleClose} />
      )}
    </>
  );
};

export default memo(JobCard);
