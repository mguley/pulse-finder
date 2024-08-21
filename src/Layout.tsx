import { FC, ReactElement, memo } from "react";
import { Box, Container } from "@mui/material";
import RoutesConfig from "./RoutesConfig";
import NavigationBar from "./domains/NavigationBar/pages/NavigationBar";

/**
 * Layout component that defines the structure of the application, including the navigation bar and main content area.
 *
 * @returns {ReactElement} The rendered Layout component.
 */
const Layout: FC = (): ReactElement => {
  return (
    <Box sx={{ flexGrow: 1 }}>
      <NavigationBar />
      <Container sx={{ mt: 4 }}>
        <RoutesConfig />
      </Container>
    </Box>
  );
};

export default memo(Layout);
