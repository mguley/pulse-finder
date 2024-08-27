import type { FC, ReactElement } from "react";
import { memo } from "react";
import { Grid } from "@mui/material";
import JobCard from "./JobCard";
import type { JobVacancy } from "../../services/jobService";

interface JobCardsProps {
  jobs: JobVacancy[];
}

/**
 * JobCards component renders a grid of JobCard components, each representing a job vacancy.
 *
 * @param {JobVacancy[]} jobs - An array of job vacancy objects to display.
 * @returns {ReactElement} The rendered component that displays a grid of job vacancies.
 */
const JobCards: FC<JobCardsProps> = ({ jobs }: JobCardsProps): ReactElement => {
  return (
    <Grid container spacing={2}>
      {jobs.map(
        (job: JobVacancy): ReactElement => (
          <Grid item xs={12} sm={6} md={4} lg={3} key={job.id}>
            <JobCard job={job} />
          </Grid>
        ),
      )}
    </Grid>
  );
};

export default memo(JobCards);
