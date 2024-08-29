import type {
  Server as HTTPServer,
  IncomingMessage,
  ServerResponse,
} from "http";
import { createServer } from "http";
import type { Socket } from "socket.io";
import { Server as SocketIOServer } from "socket.io";
import type { Tunnel } from "localtunnel";
import localtunnel from "localtunnel";
import { RecentActivityEmitter } from "../core/events/RecentActivityEmitter";
import { RecentActivityDataProvider } from "../data/RecentActivityDataProvider";
import { KeyMetricsEmitter } from "../core/events/KeyMetricsEmitter";
import { KeyMetricsDataProvider } from "../data/KeyMetricsDataProvider";
import { AESEncryption } from "../core/encryption/AESEncryption";
import type { IDataProvider } from "../core/interfaces/dataProvider/IDataProvider";
import type { IEncryption } from "../core/interfaces/encryption/IEncryption";
import type { IEventEmitter } from "../core/interfaces/event/IEventEmitter";
import type { IRecentActivity } from "../core/interfaces/recentActivity/IRecentActivity";
import type { IKeyMetrics } from "../core/interfaces/keyMetrics/IKeyMetrics";
import type { IServerManager } from "./IServerManager";

const config = {
  port: Number(process.env.PORT) || 4000,
  tunnelSubdomain: "github-io-pulse-finder",
};

/**
 * Responsible for managing the lifecycle of the HTTP and WebSocket servers, including starting, stopping,
 * handling connections, and setting up a public tunnel.
 */
export class ServerManager implements IServerManager {
  private readonly httpServer: HTTPServer;
  private readonly io: SocketIOServer;
  private tunnel?: Tunnel;

  /**
   * @param {IDataProvider<IRecentActivity>} recentActivityDataProvider - The data provider for recent activities.
   * @param {IDataProvider<IKeyMetrics>} keyMetricsDataProvider - The data provider for key metrics.
   * @param {IEncryption} encryptor - The encryption service used for securing the data.
   */
  constructor(
    private readonly recentActivityDataProvider: IDataProvider<IRecentActivity> = new RecentActivityDataProvider(),
    private readonly keyMetricsDataProvider: IDataProvider<IKeyMetrics> = new KeyMetricsDataProvider(),
    private readonly encryptor: IEncryption = new AESEncryption(),
  ) {
    this.httpServer = this.createHTTPServer();
    this.io = this.createSocketServer();
  }

  /**
   * Creates and configures the HTTP server instance.
   *
   * @returns {HTTPServer} The configured HTTP server.
   */
  private createHTTPServer(): HTTPServer {
    const requestListener = (
      request: IncomingMessage,
      response: ServerResponse,
    ) => {
      if (request.method === "OPTIONS") {
        this.handlePreflightRequest(response);
      }
    };

    return createServer(requestListener);
  }

  /**
   * Handles CORS preflight requests by responding with the appropriate headers.
   *
   * @param {ServerResponse} response - The server response object.
   */
  private handlePreflightRequest(response: ServerResponse): void {
    response.writeHead(200, {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET, POST, OPTIONS",
      "Access-Control-Allow-Headers": "Content-Type, Authorization",
    });
    response.end();
  }

  /**
   * Creates and configures the Socket.IO server instance with CORS settings.
   *
   * @returns {SocketIOServer} The configured Socket.IO server.
   */
  private createSocketServer(): SocketIOServer {
    return new SocketIOServer(this.httpServer, {
      cors: {
        origin: "*",
        methods: ["GET", "POST", "OPTIONS"],
        allowedHeaders: [
          "Content-Type",
          "bypass-tunnel-reminder",
          "Authorization",
        ],
        credentials: true,
      },
    });
  }

  /**
   * Starts the HTTP server and sets up the local tunnel for public access.
   *
   * @returns {Promise<void>} A promise that resolves when the server is successfully started and the tunnel is set up.
   */
  public async start(): Promise<void> {
    try {
      this.httpServer.listen(config.port, async () => {
        console.log(`Server is running on port ${config.port}`);
        await this.setupTunnel();
      });
      this.handleConnections();
    } catch (error) {
      console.error(`Error starting server: ${error}`);
      process.exit(1);
    }
  }

  /**
   * Stops the HTTP server and cleans up resources such as the WebSocket server and tunnel.
   *
   * @returns {Promise<void>} A promise that resolves when the server and all related resources are successfully stopped.
   */
  public async stop(): Promise<void> {
    try {
      if (this.tunnel) {
        this.tunnel.close();
      }
      this.io.close();
      this.httpServer.close((): void => {
        console.log(`Server is closed`);
      });
    } catch (error) {
      console.error(`Error stopping server: ${error}`);
    }
  }

  /**
   * Handles incoming WebSocket connections by initializing event emitters for connected clients.
   */
  public handleConnections(): void {
    this.io.on("connection", (socket: Socket) => {
      console.log(`A user connected: ${socket.id}`);
      this.initializeEmitters(socket);
    });
  }

  /**
   * Initializes the event emitters for a connected WebSocket client.
   *
   * @param {Socket} socket - The connected Socket.IO client.
   */
  private initializeEmitters(socket: Socket): void {
    try {
      if (
        this.recentActivityDataProvider &&
        this.keyMetricsDataProvider &&
        this.encryptor
      ) {
        const recentActivityEmitter: IEventEmitter = new RecentActivityEmitter(
          this.recentActivityDataProvider,
          this.encryptor,
          socket,
        );
        const keyMetricsEmitter: IEventEmitter = new KeyMetricsEmitter(
          this.keyMetricsDataProvider,
          this.encryptor,
          socket,
        );

        recentActivityEmitter.emitEvent();
        keyMetricsEmitter.emitEvent();
      }
    } catch (error) {
      console.error(`Error initializing emitters: ${error}`);
    }
  }

  /**
   * Sets up a local tunnel to expose the server publicly.
   *
   * @returns {Promise<void>} A promise that resolves when the tunnel is successfully set up.
   */
  private async setupTunnel(): Promise<void> {
    try {
      this.tunnel = await localtunnel({
        port: Number(config.port),
        subdomain: config.tunnelSubdomain,
      });
      console.log(`Server is publicly accessible via ${this.tunnel.url}`);

      this.tunnel.on("close", (): void => {
        console.log(`Tunnel closed.`);
      });
    } catch (error) {
      console.error(`Error setting up tunnel: ${error}`);
    }
  }
}
