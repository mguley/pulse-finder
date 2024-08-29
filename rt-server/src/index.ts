import { ServerManager } from "./server/ServerManager";
import { RecentActivityDataProvider } from "./data/RecentActivityDataProvider";
import { KeyMetricsDataProvider } from "./data/KeyMetricsDataProvider";
import { JobStatusDataProvider } from "./data/JobStatusDataProvider";
import { AESEncryption } from "./core/encryption/AESEncryption";

(async () => {
  // Instantiate the data providers and encryption service
  const recentActivityDataProvider = new RecentActivityDataProvider();
  const keyMetricsDataProvider = new KeyMetricsDataProvider();
  const jobStatusDataProvider = new JobStatusDataProvider();
  const encryptor = new AESEncryption();

  // Instantiate the ServerManager with the dependencies
  const serverManager = new ServerManager(
    recentActivityDataProvider,
    keyMetricsDataProvider,
    jobStatusDataProvider,
    encryptor,
  );

  await serverManager.start();

  process.on("SIGINT", async () => {
    console.log(`SIGTERM received. Stopping server...`);
    await serverManager.stop();
    process.exit(0);
  });
})();
