import type { IServerManager } from "./server/IServerManager";
import { ServerManager } from "./server/ServerManager";

// Main execution
(async () => {
  const serverManager: IServerManager = new ServerManager();
  await serverManager.start();
  serverManager.handleConnections();

  // Handle graceful shutdown
  process.on("SIGINT", async () => {
    console.log(`Received SIGINT. Shutting down...`);
    await serverManager.stop();
    process.exit(0);
  });
})();
