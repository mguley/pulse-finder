import { createServer } from "http";
import type { Socket } from "socket.io";
import { Server } from "socket.io";
import localtunnel from "localtunnel";
import { RecentActivityEmitter } from "./core/events/RecentActivityEmitter";
import { RecentActivityDataProvider } from "./data/RecentActivityDataProvider";
import { AESEncryption } from "./core/encryption/AESEncryption";
import type { IDataProvider } from "./core/interfaces/dataProvider/IDataProvider";
import type { IEncryption } from "./core/interfaces/encryption/IEncryption";
import type { IEventEmitter } from "./core/interfaces/event/IEventEmitter";
import type { IRecentActivity } from "./core/interfaces/recentActivity/IRecentActivity";

/**
 * Configures and starts the HTTP server.
 *
 * @returns {Server} - The configured Socket.IO server instance.
 */
function startHttpServer(): Server {
  const httpServer = createServer();

  const io = new Server(httpServer, {
    cors: {
      origin: "*",
      methods: ["GET", "POST"],
      allowedHeaders: ["Content-Type", "bypass-tunnel-reminder"],
    },
  });

  const PORT = process.env.PORT || 4000;
  httpServer.listen(PORT, async () => {
    console.log(`Server is running on port ${PORT}`);

    const tunnel = await localtunnel({
      port: Number(PORT),
      subdomain: "github-io-pulse-finder",
    });
    console.log(`Server is publicly accessible via ${tunnel.url}`);

    tunnel.on("close", () => console.log(`Server is closed`));
  });

  return io;
}

/**
 * Handles new WebSocket connections and initializes event emitters.
 *
 * @param {Server} io - The Socket.IO server instance.
 */
function handleConnections(io: Server): void {
  const dataProvider: IDataProvider<IRecentActivity> =
    new RecentActivityDataProvider();
  const encryptor: IEncryption = new AESEncryption();

  io.on("connection", (socket: Socket) => {
    console.log(`A user connected: ${socket.id}`);

    const activityEmitter: IEventEmitter = new RecentActivityEmitter(
      dataProvider,
      encryptor,
      socket,
    );
    activityEmitter.emitEvent();
  });
}

// Main execution
const io: Server = startHttpServer();
handleConnections(io);
