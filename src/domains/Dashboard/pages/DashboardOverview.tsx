import type { FC, ReactElement } from "react";
import { memo } from "react";
import DashboardProviders from "../components/DashboardProviders";
import KeyMetrics from "../components/keyMetrics/KeyMetrics";
import JobStatusChart from "../components/jobStatusChart/JobStatusChart";
import RecentActivityFeed from "../components/recentActivityFeed/RecentActivityFeed";
import { Box, Container } from "@mui/material";

/**
 * DashboardOverview component that provides an overview of the dashboard.
 *
 * @returns {ReactElement} The rendered DashboardOverview component.
 */
const DashboardOverview: FC = (): ReactElement => {
  return (
    <>
      <DashboardProviders>
        <Container sx={{ mt: 4, mb: 4 }}>
          <KeyMetrics />
          <Box sx={{ mt: 4 }}>
            <JobStatusChart />
          </Box>
          <Box sx={{ mt: 4 }}>
            <RecentActivityFeed />
          </Box>
        </Container>
      </DashboardProviders>
    </>
  );
};

export default memo(DashboardOverview);
