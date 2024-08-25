import { createServer } from "http";
import { Server } from "socket.io";

const httpServer = createServer();
const io = new Server(httpServer, {
  cors: {
    origin: "*",
  },
});

io.on("connection", (socket) => {
  console.log(`a user connected: ${socket.id}`);

  socket.on("disconnect", () => {
    console.log(`a user disconnected: ${socket.id}`);
  });
});

const PORT = process.env.PORT || 4000;
httpServer.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}`);
});
