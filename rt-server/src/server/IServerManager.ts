/**
 * Defines the contract for managing the lifecycle of an HTTP and WebSocket server.
 */
export interface IServerManager {
  /**
   * Starts the HTTP server and sets up the local tunnel.
   *
   * @returns {Promise<void>} A promise that resolves when the server has started.
   */
  start(): Promise<void>;

  /**
   * Handles WebSocket connections.
   */
  handleConnections(): void;

  /**
   * Closes the server and cleans up resources.
   *
   * @returns {Promise<void>} A promise that resolves when the server has stopped.
   */
  stop(): Promise<void>;
}
