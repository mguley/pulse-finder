/**
 * Represents an event handler.
 *
 * @property {string} eventName - The name of the event to listen for.
 * @property{(...args: any[]) => void} handler - The event handler function to invoke when the event is emitted.
 */
export interface EventHandler {
  eventName: string;
  handler: (...args: any[]) => void;
}

/**
 * Represents a web socket connection.
 */
export interface WebSocket {
  /**
   * Connects to the WebSocket server.
   */
  connect(): void;

  /**
   * Registers a new event listener on the socket.
   *
   * @param {EventHandler} eventHandler - The event handler object containing the event name and handler function.
   */
  on(eventHandler: EventHandler): void;

  /**
   * Removes a specific event listener from the socket.
   *
   * @param {EventHandler} eventHandler - The event handler object containing the event name and handler function.
   */
  off(eventHandler: EventHandler): void;

  /**
   * Disconnects from the WebSocket server.
   */
  disconnect(): void;

  /**
   * Checks if the WebSocket is connected.
   *
   * @returns {boolean} - Returns true if the socket is connected, false otherwise.
   */
  isConnected(): boolean;
}
