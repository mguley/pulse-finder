import { FC, ReactElement, memo } from "react";
import { Route, Routes } from "react-router-dom";
import DashboardOverview from "./domains/Dashboard/pages/DashboardOverview";
import JobsManagement from "./domains/Jobs/pages/JobsManagement";
import WorkerMonitoring from "./domains/Workers/pages/WorkerMonitoring";

/**
 * RoutesConfig component that defines all the routes for the application.
 *
 * @returns {ReactElement} The rendered RoutesConfig component.
 */
const RoutesConfig: FC = (): ReactElement => {
  return (
    <Routes>
      <Route path="/pulse-finder/" element={<DashboardOverview />} />
      <Route
        path="/pulse-finder/jobs-management"
        element={<JobsManagement />}
      />
      <Route
        path="/pulse-finder/worker-monitoring"
        element={<WorkerMonitoring />}
      />
    </Routes>
  );
};

export default memo(RoutesConfig);
