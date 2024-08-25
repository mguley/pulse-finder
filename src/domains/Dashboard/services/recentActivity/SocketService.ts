import { io, Socket } from "socket.io-client";
import { RecentActivitySocket } from "./types";

/**
 * SocketService class to manage the WebSocket connection.
 *
 * Responsible for establishing, managing, and terminating a WebSocket connection to a specified server URL.
 * It provides methods to handle events, check connection status, and disconnect when necessary.
 */
export class SocketService implements RecentActivitySocket {
  private socket: Socket | null = null;
  private readonly serverUrl: string;

  /**
   * Creates an instance of SocketService.
   *
   * @param {string} serverUrl - The URL of the WebSocket server to connect to.
   */
  constructor(serverUrl: string) {
    this.serverUrl = serverUrl;
  }

  /**
   * Connects to the WebSocket server.
   *
   * Establishes a WebSocket connection to the server specified in the constructor.
   * It also sets up listeners for connection and disconnection events.
   */
  public connect(): void {
    if (!this.socket) {
      this.socket = io(this.serverUrl);

      this.socket.on("connect", () => {
        console.log(
          `Connected to ${this.serverUrl} with id: ${this.socket!.id}`,
        );
      });

      this.socket.on("disconnect", () => {
        console.log(
          `Disconnected from ${this.serverUrl} with id: ${this.socket!.id}`,
        );
      });
    }
  }

  /**
   * Disconnects from the WebSocket server.
   *
   * Closes the WebSocket connection and cleans up the socket instance.
   */
  disconnect(): void {
    if (this.socket) {
      this.socket.disconnect();
      this.socket = null;
    }
  }

  /**
   * Checks if the WebSocket is connected.
   *
   * @returns {boolean} - Returns true if the socket is connected, false otherwise.
   */
  isConnected(): boolean {
    return this.socket !== null && this.socket.connected;
  }

  /**
   * Removes a specific event listener from the socket.
   *
   * @param {string} eventName - The name of the event to stop listening for.
   * @param {(...args: any[]) => void} handler - The event handler function to remove.
   */
  off(eventName: string, handler: (...args: any[]) => void): void {
    if (this.socket) {
      this.socket.off(eventName, handler);
    }
  }

  /**
   * Registers a new event listener on the socket.
   *
   * @param {string} eventName - The name of the event to listen for.
   * @param {(...args: any[]) => void} handler - The event handler function to invoke when the event is emitted.
   */
  on(eventName: string, handler: (...args: any[]) => void): void {
    if (this.socket) {
      this.socket.on(eventName, handler);
    }
  }
}
