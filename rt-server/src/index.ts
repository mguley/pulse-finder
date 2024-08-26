import { createServer } from "http";
import { Server } from "socket.io";
import localtunnel from "localtunnel";
import {
  ActivityEmitter,
  RandomActivityGenerator,
} from "./emitter/ActivityEmitter";
import { mockActivities } from "./mock/mockData";

const httpServer = createServer();
const io = new Server(httpServer, {
  cors: {
    origin: "*",
  },
});

const activityGenerator = new RandomActivityGenerator(mockActivities);
const activityEmitter = new ActivityEmitter(io, activityGenerator, {
  intervalMs: 5000,
});

io.on("connection", (socket) => {
  console.log(`A user connected: ${socket.id}`);
  activityEmitter.startEmitting(socket);
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
