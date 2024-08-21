import { FC, ReactElement, memo } from "react";
import Typography from "@mui/material/Typography";

const DashboardOverview: FC = (): ReactElement => {
  return <Typography variant="h4">Dashboard Overview</Typography>;
};

export default memo(DashboardOverview);
