import { createServer } from "http";
import { Server } from "socket.io";
import localtunnel from "localtunnel";
import {
  ActivityEmitter,
  RandomActivityGenerator,
} from "./emitter/ActivityEmitter";
import { mockActivities } from "./mock/mockData";

/**
 * Creates and configures the HTTP server and Socket.IO server.
 */
const httpServer = createServer();

/**
 * Initializes the Socket.IO server with CORS configuration.
 *
 * CORS (Cross-Origin Resource Sharing) is configured to allow all origins,
 * and permits GET and POST methods. Additionally, a custom header
 * 'bypass-tunnel-reminder' is allowed to bypass the localtunnel reminder page.
 */
const io = new Server(httpServer, {
  cors: {
    origin: "*",
    methods: ["GET", "POST"],
    allowedHeaders: ["Content-Type", "bypass-tunnel-reminder"],
  },
});

/**
 * Creates an instance of RandomActivityGenerator with mock data.
 * This instance will be used to generate random activities to emit.
 */
const activityGenerator = new RandomActivityGenerator(mockActivities);

/**
 * Creates an instance of ActivityEmitter to manage the emission of activities.
 * Activities are emitted every 5 seconds by default.
 */
const activityEmitter = new ActivityEmitter(io, activityGenerator, {
  intervalMs: 5000,
});

/**
 * Handles new WebSocket connections.
 *
 * When a new client connects, the server starts emitting activities to that client.
 * The ActivityEmitter is responsible for managing the emission process.
 */
io.on("connection", (socket) => {
  console.log(`A user connected: ${socket.id}`);
  activityEmitter.startEmitting(socket);
});

/**
 * Starts the HTTP server and exposes it via a public URL using localtunnel.
 *
 * The server listens on a specified port (default 4000).
 * It also sets up a tunnel using localtunnel to make the server publicly accessible.
 * The subdomain 'github-io-pulse-finder' is used for the public URL.
 */
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
