import type { FC, ReactElement } from "react";
import { memo } from "react";
import { BrowserRouter as Router } from "react-router-dom";
import Layout from "./Layout";
import { NavigationBarProvider } from "./domains/NavigationBar/context/NavigationBarProvider";

/**
 * App component that serves as the root component of the application.
 *
 * @returns {ReactElement} The rendered App component.
 */
const App: FC = (): ReactElement => {
  return (
    <NavigationBarProvider>
      <Router>
        <Layout />
      </Router>
    </NavigationBarProvider>
  );
};

export default memo(App);
