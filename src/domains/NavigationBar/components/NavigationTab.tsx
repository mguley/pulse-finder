import type { FC, ReactElement } from "react";
import { memo } from "react";
import { Tab } from "@mui/material";

interface NavigationTabProps {
  key: number;
  label: string;
  isActive: boolean;
  onClick: () => void;
}

/**
 * NavigationTab component that renders a single tab with custom styles based on its active state.
 *
 * @param {number} key - A unique key for the tab, usually the index of the tab.
 * @param {string} label - The text label of the tab.
 * @param {boolean} isActive - Whether the tab is currently active.
 * @param {Function} onClick - The function to call when the tab is clicked.
 * @returns {ReactElement} - The rendered NavigationTab component.
 */
const NavigationTab: FC<NavigationTabProps> = ({
  key,
  label,
  isActive,
  onClick,
}: NavigationTabProps): ReactElement => {
  return (
    <Tab
      key={key}
      label={label}
      onClick={onClick}
      sx={{
        color: isActive ? "white" : "rgba(255, 255, 255, 0.7)",
        fontWeight: isActive ? "bold" : "normal",
      }}
    />
  );
};

export default memo(NavigationTab);
