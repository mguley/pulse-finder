import type { FC, ReactElement } from "react";
import { memo } from "react";
import { Route, Routes } from "react-router-dom";
import JobsManagement from "./domains/Jobs/pages/JobsManagement";

/**
 * RoutesConfig component that defines all the routes for the application.
 *
 * @returns {ReactElement} The rendered RoutesConfig component.
 */
const RoutesConfig: FC = (): ReactElement => {
  return (
    <Routes>
      <Route path="/pulse-finder/" element={<JobsManagement />} />
    </Routes>
  );
};

export default memo(RoutesConfig);
