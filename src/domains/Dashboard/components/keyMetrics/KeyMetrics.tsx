import type { FC, ReactElement } from "react";
import { memo } from "react";
import { Grid, Card, CardContent, Typography } from "@mui/material";
import { useKeyMetrics } from "../../context/KeyMetricsContext";
import type { KeyMetricService as IKeyMetric } from "../../services/keyMetricService";
import LoadingSpinner from "../../../../shared/components/LoadingSpinner";

/**
 * KeyMetrics component that displays key metrics cards in the Dashboard overview.
 *
 * @returns {ReactElement} The rendered KeyMetrics component.
 */
const KeyMetrics: FC = (): ReactElement => {
  const { keyMetrics, loading, error } = useKeyMetrics();

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

  if (!keyMetrics || keyMetrics.length === 0) {
    return <Typography variant="h6">No metrics available.</Typography>;
  }

  return (
    <>
      <Typography variant="h6" align="center">
        {" "}
        Key Metrics Cards
      </Typography>
      <Grid container spacing={3}>
        {keyMetrics.map(({ title, value }: IKeyMetric, index: number) => (
          <Grid item xs={12} sm={4} key={index}>
            <Card sx={{ border: "1 px solid #ccc", boxShadow: 3 }}>
              <CardContent>
                <Typography variant="h6">{title}</Typography>
                <Typography variant="h5">{value}</Typography>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </>
  );
};

export default memo(KeyMetrics);
