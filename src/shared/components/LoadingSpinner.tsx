import { FC, ReactElement, memo } from "react";
import { CircularProgress, Box } from "@mui/material";

/**
 * LoadingSpinner component that displays a circular progress indicator.
 *
 * @returns {ReactElement} The rendered LoadingSpinner component.
 */
const LoadingSpinner: FC = (): ReactElement => {
  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "100px",
      }}
    >
      <CircularProgress />
    </Box>
  );
};

export default memo(LoadingSpinner);
