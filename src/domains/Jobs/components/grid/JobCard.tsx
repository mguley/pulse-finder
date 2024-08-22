import { FC, ReactElement, memo } from "react";
import { Card, CardContent, Typography } from "@mui/material";
import { JobVacancy } from "../../services/jobService";

interface JobCardProps {
  job: JobVacancy;
}

/**
 * JobCard component is a presentational component that displays the details of a single job vacancy.
 *
 * @param {JobVacancy} job - The job vacancy object containing details to display.
 * @returns {ReactElement} The rendered component that displays a job vacancy.
 */
const JobCard: FC<JobCardProps> = ({
  job: { title, company },
}: JobCardProps): ReactElement => {
  return (
    <Card variant="outlined">
      <CardContent>
        <Typography variant="h6" component="div">
          {title}
        </Typography>
        <Typography color="text.secondary">{company}</Typography>
      </CardContent>
    </Card>
  );
};

export default memo(JobCard);
