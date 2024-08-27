import type { FC, ReactElement } from "react";
import { memo } from "react";
import { Tabs, Box } from "@mui/material";
import { useLocation, useNavigate } from "react-router-dom";
import type { NavigationItem } from "../context/NavigationBarProvider";
import { useNavigationBar } from "../context/NavigationBarProvider";
import NavigationTab from "./NavigationTab";

/**
 * NavigationTabs component that renders a collection of navigation tabs based on the provided navigation items.
 * It handles the logic for determining the active tab and changing routes when a tab is selected.
 *
 * @returns {ReactElement} The rendered NavigationTabs component.
 */
const NavigationTabs: FC = (): ReactElement => {
  const { navItems } = useNavigationBar();
  const location = useLocation();
  const navigate = useNavigate();

  /**
   * Determines current active tab based on the current route.
   *
   * @type {number} - The index of the currently active tab.
   */
  const currentTab = navItems.findIndex(
    (item: NavigationItem) => item.route === location.pathname,
  );

  /**
   * Handles tab change events. When a tab is selected, it navigates to the corresponding route.
   *
   * @param {number} newValue - The index of the newly selected tab.
   */
  const handleTabChange = (newValue: number): void => {
    navigate(navItems[newValue].route);
  };

  return (
    <Box sx={{ display: "flex", flexGrow: 1, justifyContent: "center" }}>
      <Tabs
        value={currentTab}
        centered
        textColor="inherit"
        TabIndicatorProps={{ style: { backgroundColor: "white" } }}
      >
        {navItems.map((item: NavigationItem, index: number) => (
          <NavigationTab
            key={index}
            label={item.label}
            isActive={currentTab === index}
            onClick={() => handleTabChange(index)}
          />
        ))}
      </Tabs>
    </Box>
  );
};

export default memo(NavigationTabs);
