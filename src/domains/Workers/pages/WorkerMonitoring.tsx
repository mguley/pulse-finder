import { FC, ReactElement, memo } from "react";
import Typography from "@mui/material/Typography";

const WorkerMonitoring: FC = (): ReactElement => {
  return <Typography variant="h4">Workers Management</Typography>;
};

export default memo(WorkerMonitoring);
