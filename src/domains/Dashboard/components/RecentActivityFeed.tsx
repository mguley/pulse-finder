import { FC, ReactElement, memo } from "react";
import { Card, CardContent, Typography, Box, Grid } from "@mui/material";
import { useRecentActivityFeed } from "../context/RecentActivityFeedContext";
import { RecentActivity } from "../services/recentActivity/types";
import LoadingSpinner from "../../../shared/components/LoadingSpinner";

/**
 * RecentActivityFeed component that displays a list of recent job activities.
 *
 * @returns {ReactElement} The rendered RecentActivityFeed component.
 */
const RecentActivityFeed: FC = (): ReactElement => {
  const { recentActivities, loading, error } = useRecentActivityFeed();

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

  if (!recentActivities || recentActivities.length === 0) {
    return <Typography variant="h6">No recent activity available.</Typography>;
  }

  return (
    <Card sx={{ border: "1px solid #ccc", boxShadow: 3 }}>
      <CardContent>
        <Typography variant="h6" align="center">
          Recent Activity
        </Typography>
        <Box sx={{ mt: 2 }}>
          {/* Column headers */}
          <Grid container spacing={2} sx={{ fontWeight: "bold", mb: 1 }}>
            <Grid item xs={3}>
              <Typography>Job ID</Typography>
            </Grid>
            <Grid item xs={3}>
              <Typography>Status</Typography>
            </Grid>
            <Grid item xs={3}>
              <Typography>Start Time</Typography>
            </Grid>
            <Grid item xs={3}>
              <Typography>End Time</Typography>
            </Grid>
          </Grid>
          {/* Data Rows */}
          {recentActivities.map(
            (
              { jobId, status, startTime, endTime }: RecentActivity,
              index: number,
            ) => (
              <Grid container spacing={2} key={index}>
                <Grid item xs={3}>
                  <Typography>{jobId}</Typography>
                </Grid>
                <Grid item xs={3}>
                  <Typography>{status}</Typography>
                </Grid>
                <Grid item xs={3}>
                  <Typography>{startTime}</Typography>
                </Grid>
                <Grid item xs={3}>
                  <Typography>{endTime}</Typography>
                </Grid>
              </Grid>
            ),
          )}
        </Box>
      </CardContent>
    </Card>
  );
};

export default memo(RecentActivityFeed);
