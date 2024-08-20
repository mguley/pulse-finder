import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "../src/App";

/**
 * The entry point for rendering the React application in the DOM.
 */
const container = document.getElementById("root");
if (!container) {
  throw new Error("Could not find container");
}

const root = createRoot(container);

root.render(
  <StrictMode>
    <App />
  </StrictMode>,
);
