import type { Server as HTTPServer } from "http";
import type { Socket } from "socket.io";
import { Server as SocketIOServer } from "socket.io";
import type { IDataProvider } from "../core/interfaces/dataProvider/IDataProvider";
import type { IEncryption } from "../core/interfaces/encryption/IEncryption";
import type { IRecentActivity } from "../core/interfaces/recentActivity/IRecentActivity";
import type { IKeyMetrics } from "../core/interfaces/keyMetrics/IKeyMetrics";
import { EventEmitterManager } from "../core/events/EventEmitterManager";
import type { IJobStatus } from "../core/interfaces/jobStatus/IJobStatus";

/**
 * Responsible for managing the lifecycle of a WebSocket server.
 */
export class WebSocketServerManager {
  private readonly io: SocketIOServer;

  /**
   * @param {HTTPServer} httpServer - The HTTP server instance to which the WebSocket server will be attached.
   */
  constructor(httpServer: HTTPServer) {
    this.io = this.createSocketServer(httpServer);
  }

  /**
   * Creates and configures the WebSocket server.
   *
   * @param {HTTPServer} httpServer - The HTTP server instance to attach the WebSocket server to.
   * @returns {SocketIOServer} - The configured WebSocket server instance.
   */
  private createSocketServer(httpServer: HTTPServer): SocketIOServer {
    return new SocketIOServer(httpServer, {
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
   * Handles incoming WebSocket connections.
   *
   * @param {IDataProvider<IRecentActivity>} recentActivityDataProvider - The data provider for recent activity data.
   * @param {IDataProvider<IKeyMetrics>} keyMetricsDataProvider - The data provider for key metrics data.
   * @param {IDataProvider<IJobStatus>} jobStatusDataProvider - The data provider for job status data.
   * @param {IEncryption} encryptor - The encryption service used to secure data before it is sent to clients.
   */
  public handleConnections(
    recentActivityDataProvider: IDataProvider<IRecentActivity>,
    keyMetricsDataProvider: IDataProvider<IKeyMetrics>,
    jobStatusDataProvider: IDataProvider<IJobStatus>,
    encryptor: IEncryption,
  ): void {
    this.io.on("connection", (socket: Socket) => {
      console.log(`A user connected: ${socket.id}`);
      const eventEmitterFactory = new EventEmitterManager(
        recentActivityDataProvider,
        keyMetricsDataProvider,
        jobStatusDataProvider,
        encryptor,
        socket,
      );
      eventEmitterFactory.emit();
    });
  }

  /**
   * Closes the WebSocket server.
   */
  public close(): void {
    this.io.close((): void => {
      console.log(`WebSocket Server is closed`);
    });
  }
}
