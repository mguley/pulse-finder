import { FC, ReactElement, memo } from "react";
import { KeyMetricsProvider } from "../context/KeyMetricsContext";
import { JobStatusChartProvider } from "../context/JobStatusChartContext";
import KeyMetrics from "../components/KeyMetrics";
import JobStatusChart from "../components/JobStatusChart";
import { Box, Container } from "@mui/material";

/**
 * DashboardOverview component that provides an overview of the dashboard.
 *
 * @returns {ReactElement} The rendered DashboardOverview component.
 */
const DashboardOverview: FC = (): ReactElement => {
  return (
    <>
      <Container sx={{ mt: 4, mb: 4 }}>
        <KeyMetricsProvider>
          <KeyMetrics />
        </KeyMetricsProvider>
        <Box sx={{ mt: 4 }}>
          <JobStatusChartProvider>
            <JobStatusChart />
          </JobStatusChartProvider>
        </Box>
      </Container>
    </>
  );
};

export default memo(DashboardOverview);
