import { FC, ReactElement, memo } from "react";

/**
 * App component that serves as the root component of the application.
 *
 * @returns {ReactElement} The rendered App component.
 */
const App: FC = (): ReactElement => {
    return <h1>Hello</h1>
};

export default memo(App);
