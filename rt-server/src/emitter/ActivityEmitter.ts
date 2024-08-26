import { Server, Socket } from "socket.io";
import { RecentActivity } from "../types";

/**
 * Represents a generic activity generator.
 *
 * Implementations of this interface should provide logic for generating activities
 * that can be emitted via a WebSocket connection.
 */
export interface ActivityGenerator {
  /**
   * Generates a recent activity.
   *
   * @returns {RecentActivity} - A randomly generated or predefined recent activity.
   */
  generateActivity(): RecentActivity;
}

/**
 * Generates random recent activities from a predefined list of activities.
 */
export class RandomActivityGenerator implements ActivityGenerator {
  private readonly activities: RecentActivity[] = [];

  /**
   * @param {RecentActivity[]} activities - An array of recent activities to randomly select from.
   * @throws Will throw an error if the provided activity list is empty.
   */
  constructor(activities: RecentActivity[]) {
    if (activities.length === 0) {
      throw new Error(`Activity list cannot be empty`);
    }
    this.activities = activities;
  }

  /**
   * Generates a random activity from the provided list.
   *
   * @returns {RecentActivity} - A randomly selected recent activity.
   */
  public generateActivity(): RecentActivity {
    const randomIndex = Math.floor(Math.random() * this.activities.length);
    return this.activities[randomIndex];
  }
}

/**
 * Configuration options.
 *
 * @property {number} intervalMs - The interval, in milliseconds, between each activity emission.
 */
export interface ActivityEmitterConfig {
  intervalMs?: number;
}

/**
 * Responsible for emitting recent activities to connected clients via a WebSocket connection.
 * Uses an instance of ActivityGenerator to generate the activities and manages the timing of
 * these emissions based on a configurable interval.
 */
export class ActivityEmitter {
  private io: Server;
  private activityGenerator: ActivityGenerator;
  private readonly intervalMs: number;

  /**
   * @param {Server} io - The Socket.IO server instance used to manage WebSocket connections.
   * @param {ActivityGenerator} activityGenerator - An instance of ActivityGenerator used to generate activities.
   * @param {ActivityEmitterConfig} config - Optional configuration for the emission interval.
   */
  constructor(
    io: Server,
    activityGenerator: ActivityGenerator,
    config: ActivityEmitterConfig = {},
  ) {
    this.io = io;
    this.activityGenerator = activityGenerator;
    this.intervalMs = config.intervalMs || 5000;
  }

  /**
   * Starts emitting activities to a connected client.
   * It sends an initial activity immediately upon a client's connection, and continues to emit activities
   * at regular intervals until the client disconnects.
   *
   * @param {Socket} socket - The Socket.IO socket instance representing the connected client.
   */
  public startEmitting(socket: Socket): void {
    const sendRandomActivity = (): void => {
      const randomActivity = this.activityGenerator.generateActivity();
      socket.emit("newActivity", randomActivity);
      console.log(
        `Sent activity to ${socket.id}: ${JSON.stringify(randomActivity)}`,
      );
    };

    sendRandomActivity();
    const intervalId = setInterval(sendRandomActivity, this.intervalMs);

    socket.on("disconnect", () => {
      clearInterval(intervalId);
      console.log(`A user disconnected: ${socket.id}`);
    });
  }
}
