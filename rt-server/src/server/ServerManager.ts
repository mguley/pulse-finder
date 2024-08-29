import type { IDataProvider } from "../core/interfaces/dataProvider/IDataProvider";
import type { IEncryption } from "../core/interfaces/encryption/IEncryption";
import type { IRecentActivity } from "../core/interfaces/recentActivity/IRecentActivity";
import type { IKeyMetrics } from "../core/interfaces/keyMetrics/IKeyMetrics";
import { HTTPServerManager } from "./HTTPServerManager";
import { WebSocketServerManager } from "./WebSocketServerManager";
import { TunnelManager } from "./TunnelManager";
import type { IJobStatus } from "../core/interfaces/jobStatus/IJobStatus";

const config = {
  port: Number(process.env.PORT) || 4000,
  tunnelSubdomain: "github-io-pulse-finder",
};

/**
 * Responsible for managing the lifecycle of the HTTP and WebSocket servers, including starting, stopping,
 * handling connections, and setting up a public tunnel.
 */
export class ServerManager {
  private readonly httpServerManager: HTTPServerManager;
  private readonly webSocketServerManager: WebSocketServerManager;
  private readonly tunnelManager: TunnelManager;

  /**
   * @param {IDataProvider<IRecentActivity>} recentActivityDataProvider - The data provider for recent activity data.
   * @param {IDataProvider<IKeyMetrics>} keyMetricsDataProvider - The data provider for key metrics data.
   * @param {IDataProvider<IJobStatus>} jobStatusDataProvider - The data provider for job status data.
   * @param {IEncryption} encryptor - The encryption service used to secure data before transmission.
   */
  constructor(
    private readonly recentActivityDataProvider: IDataProvider<IRecentActivity>,
    private readonly keyMetricsDataProvider: IDataProvider<IKeyMetrics>,
    private readonly jobStatusDataProvider: IDataProvider<IJobStatus>,
    private readonly encryptor: IEncryption,
  ) {
    this.httpServerManager = new HTTPServerManager();
    this.webSocketServerManager = new WebSocketServerManager(
      this.httpServerManager.getServer(),
    );
    this.tunnelManager = new TunnelManager(config.port, config.tunnelSubdomain);
  }

  /**
   * Starts the server components, including the HTTP server, WebSocket server, and public tunnel.
   *
   * @returns {Promise<void>} A promise that resolves when the server has started successfully.
   */
  public async start(): Promise<void> {
    try {
      // Start the HTTP server
      this.httpServerManager.start();

      // Handle WebSocket connections
      this.webSocketServerManager.handleConnections(
        this.recentActivityDataProvider,
        this.keyMetricsDataProvider,
        this.jobStatusDataProvider,
        this.encryptor,
      );

      // Set up the public tunnel
      await this.tunnelManager.setupTunnel();
    } catch (e) {
      console.error(`Error starting server: ${e}`);
      process.exit(1);
    }
  }

  /**
   * Stops the server components, including HTTP server, WebSocket server, and public tunnel.
   *
   * @returns {Promise<void>} A promise that resolves when the server has stopped successfully.
   */
  public async stop(): Promise<void> {
    try {
      await this.tunnelManager.closeTunnel();
      this.webSocketServerManager.close();
      this.httpServerManager.stop();
    } catch (e) {
      console.error(`Error stopping server: ${e}`);
    }
  }
}
