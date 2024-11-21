import type { FC, ReactElement } from "react";
import { memo } from "react";
import { Modal, Box, Typography } from "@mui/material";
import type { JobVacancy } from "../../services/jobs/response";

interface JobDetailsModalProps {
  job: JobVacancy;
  open: boolean;
  onClose: () => void;
}

/**
 * JobDetailsModal component renders detailed information about a job vacancy in a modal window.
 *
 * @param {JobVacancy} job - The job vacancy object containing details to display.
 * @param {boolean} open - A flag indicating whether the modal is open or closed.
 * @param {function} onClose - A function to handle closing the modal.
 * @returns {ReactElement} The rendered component that displays job details in a modal.
 */
const JobDetailsModal: FC<JobDetailsModalProps> = ({
  job,
  open,
  onClose,
}: JobDetailsModalProps): ReactElement => {
  return (
    <Modal
      open={open}
      onClose={onClose}
      aria-labelledby="job-title"
      aria-describedby="job-description"
      slotProps={{
        backdrop: {
          sx: {
            backgroundColor: "rgba(0, 0, 0, 0.1)", // Lightly darkened background
            transition: "all 1s ease-in-out", // Smooth transition
          },
        },
      }}
    >
      <Box
        sx={{
          position: "absolute",
          top: "50%",
          left: "50%",
          transform: "translate(-50%, -50%)",
          width: 400,
          bgcolor: "background.paper",
          boxShadow: 5,
          p: 4,
        }}
      >
        <Typography id="job-title" variant="h6" component="h2">
          {job.title}
        </Typography>
        <Typography sx={{ mt: 2 }}>
          Job Description: {job.description}
        </Typography>
        <Typography sx={{ mt: 2 }}>Posted At: {job.posted_at}</Typography>
        <Typography sx={{ mt: 2 }}>Location: {job.location}</Typography>
      </Box>
    </Modal>
  );
};

export default memo(JobDetailsModal);
