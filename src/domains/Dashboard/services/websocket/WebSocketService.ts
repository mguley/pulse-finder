import { io, Socket } from "socket.io-client";
import { WebSocket, EventHandler } from "./types";

/**
 * Manages a WebSocket connection.
 * Responsible for establishing, managing, and terminating a WebSocket connection to a specified server URL.
 * It provides methods to handle events, check connection status, and disconnect when necessary.
 */
export class SocketService implements WebSocket {
  private socket: Socket | null = null;
  private readonly serverUrl: string = "";

  /**
   * Creates an instance of SocketService.
   *
   * @param {string} serverUrl - The URL of the WebSocket server to connect to.
   */
  constructor(serverUrl?: string) {
    this.serverUrl = serverUrl ?? "https://github-io-pulse-finder.loca.lt";
  }

  /**
   * Connects to the WebSocket server.
   */
  public connect(): void {
    if (!this.socket) {
      this.socket = io(this.serverUrl, {
        extraHeaders: {
          "bypass-tunnel-reminder": "true", // Bypass the tunnel reminder
        },
      });
    }
  }

  /**
   * Disconnects from the WebSocket server.
   * It closes the WebSocket connection and cleans up the socket instance.
   */
  public disconnect(): void {
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
  public isConnected(): boolean {
    return this.socket !== null && this.socket.connected;
  }

  /**
   * Removes a specific event listener from the socket.
   *
   * @param {EventHandler} eventHandler - The event handler object containing the event name and handler function.
   */
  public off(eventHandler: EventHandler): void {
    if (this.socket) {
      this.socket.off(eventHandler.eventName, eventHandler.handler);
    }
  }

  /**
   * Registers a new event listener on the socket.
   *
   * @param {EventHandler} eventHandler - The event handler object containing the event name and handler function.
   */
  public on(eventHandler: EventHandler): void {
    if (this.socket) {
      this.socket.on(eventHandler.eventName, eventHandler.handler);
    }
  }
}
