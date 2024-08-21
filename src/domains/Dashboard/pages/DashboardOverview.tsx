import { FC, ReactElement, memo } from "react";
import { KeyMetricsProvider } from "../context/KeyMetricsContext";
import KeyMetrics from "../components/KeyMetrics";

/**
 * DashboardOverview component that provides an overview of the dashboard.
 *
 * @returns {ReactElement} The rendered DashboardOverview component.
 */
const DashboardOverview: FC = (): ReactElement => {
  return (
    <>
      <KeyMetricsProvider>
        <KeyMetrics />
      </KeyMetricsProvider>
    </>
  );
};

export default memo(DashboardOverview);
