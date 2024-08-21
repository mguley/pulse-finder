import { FC, ReactElement, memo } from "react";
import { AppBar } from "@mui/material";
import NavigationTabs from "../components/NavigationTabs";

/**
 * NavigationBar component that renders the top navigation bar.
 * It uses the NavigationTabs component to manage and display the tabs.
 *
 * @returns {ReactElement} The rendered NavigationBar component.
 */
const NavigationBar: FC = (): ReactElement => {
  return (
    <AppBar position="static">
      <NavigationTabs />
    </AppBar>
  );
};

export default memo(NavigationBar);
