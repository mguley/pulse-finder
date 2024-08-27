import type { FC, ReactElement, ReactNode } from "react";
import { memo } from "react";
import { KeyMetricsProvider } from "../context/KeyMetricsContext";
import { JobStatusChartProvider } from "../context/JobStatusChartContext";
import { RecentActivityFeedProvider } from "../context/RecentActivityFeedContext";

interface DashboardProvidersProps {
  children: ReactNode;
}

/**
 * DashboardProviders component consolidates multiple context providers for the Dashboard.
 * It serves as a wrapper that provides all necessary contexts to its children components, ensuring
 * that they have access to the appropriate state and data.
 *
 * Contexts Provided:
 * - KeyMetricsContext: Manages the state and data for key metrics displayed on the dashboard.
 * - JobStatusChartContext: Manages the state and data for the job status chart.
 * - RecentActivityFeedContext: Manages the state and data for the recent activity.
 *
 * @param {DashboardProvidersProps} props - The props for the DashboardProviders component.
 * @returns {ReactElement} The rendered component that provides the required contexts to its children.
 */
const DashboardProviders: FC<DashboardProvidersProps> = ({
  children,
}: DashboardProvidersProps): ReactElement => {
  return (
    <KeyMetricsProvider>
      <JobStatusChartProvider>
        <RecentActivityFeedProvider>{children}</RecentActivityFeedProvider>
      </JobStatusChartProvider>
    </KeyMetricsProvider>
  );
};

export default memo(DashboardProviders);
