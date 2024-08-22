import { FC, ReactElement, memo } from "react";
import { PieChart } from "@mui/x-charts/PieChart";
import { Card, CardContent, Typography } from "@mui/material";
import { useJobStatusChart } from "../context/JobStatusChartContext";
import { JobStatusService as IJobStatus } from "../services/jobStatusService";
import LoadingSpinner from "../../../shared/components/LoadingSpinner";

/**
 * JobStatusChart component that displays a pie chart of jobs statuses.
 *
 * @returns {ReactElement} The rendered JobStatusChart component.
 */
const JobStatusChart: FC = (): ReactElement => {
  const { jobStatusData, loading, error } = useJobStatusChart();

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

  if (!jobStatusData || jobStatusData.length === 0) {
    return <Typography variant="h6">No job status data available.</Typography>;
  }

  const data = jobStatusData.map(
    ({ status, count }: IJobStatus, index: number) => ({
      id: index,
      value: count,
      label: status,
    }),
  );

  return (
    <Card sx={{ border: "1px solid #ccc", boxShadow: 3 }}>
      <CardContent>
        <Typography variant="h6" align="center">
          Job Statuses
        </Typography>
        <PieChart
          width={400}
          height={200}
          series={[
            {
              data,
            },
          ]}
        />
      </CardContent>
    </Card>
  );
};

export default memo(JobStatusChart);
