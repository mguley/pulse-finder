import { createServer } from "http";
import { Server } from "socket.io";
import { mockActivities } from "./mockData";

const httpServer = createServer();
const io = new Server(httpServer, {
  cors: {
    origin: "*",
  },
});

io.on("connection", (socket) => {
  console.log(`A user connected: ${socket.id}`);

  // Function to emit a random activity every 5 seconds
  const sendRandomActivity = (): void => {
    const randomIndex = Math.floor(Math.random() * mockActivities.length);
    const randomActivity = mockActivities[randomIndex];

    socket.emit("newActivity", randomActivity);
    console.log(
      `Sent activity to ${socket.id}: ${JSON.stringify(randomActivity)}`,
    );
  };

  // Send an activity immediately upon connection
  sendRandomActivity();

  // Set interval to send an activity every 5 seconds
  const intervalId = setInterval(sendRandomActivity, 5000);

  // Clear the interval when the user disconnects
  socket.on("disconnect", () => {
    clearInterval(intervalId);
    console.log(`A user disconnected: ${socket.id}`);
  });
});

const PORT = process.env.PORT || 4000;
httpServer.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}`);
});
