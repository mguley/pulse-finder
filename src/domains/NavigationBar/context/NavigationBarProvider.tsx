import type { FC, ReactNode, ReactElement } from "react";
import { createContext, useContext, useState } from "react";

export interface NavigationItem {
  label: string;
  route: string;
}

interface NavigationBarContextType {
  navItems: NavigationItem[];
}

interface NavigationBarProviderProps {
  children: ReactNode;
}

const NavigationBarContext = createContext<
  NavigationBarContextType | undefined
>(undefined);

/**
 * NavigationBarProvider component that provides the navigation items to its children via context.
 *
 * @param {NavigationBarProviderProps} props - The properties for the NavigationBarProvider component.
 * @returns {ReactElement} The rendered NavigationBarProvider component.
 */
export const NavigationBarProvider: FC<NavigationBarProviderProps> = ({
  children,
}: NavigationBarProviderProps): ReactElement => {
  const data: NavigationItem[] = [
    { label: "Dashboard", route: "/pulse-finder/" },
    { label: "Jobs Management", route: "/pulse-finder/jobs-management" },
    { label: "Worker Monitoring", route: "/pulse-finder/worker-monitoring" },
  ];

  const [navItems] = useState<NavigationItem[]>(data);

  return (
    <NavigationBarContext.Provider value={{ navItems }}>
      {children}
    </NavigationBarContext.Provider>
  );
};

/**
 * useNavigationBar hook that provides access to the navigation bar context.
 *
 * @returns {NavigationBarContextType} The context value containing the navigation items.
 */
export const useNavigationBar = (): NavigationBarContextType => {
  const context: NavigationBarContextType | undefined =
    useContext(NavigationBarContext);
  if (!context) {
    throw new Error(
      "useNavigationBar must be used within a NavigationBarProvider",
    );
  }
  return context;
};
