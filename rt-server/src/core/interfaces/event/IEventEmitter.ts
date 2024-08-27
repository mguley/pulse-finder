/**
 * Defines the contract for event emitter.
 */
export interface IEventEmitter {
  /**
   * Emits an event.
   *
   * This method should contain the logic necessary to broadcast or trigger an event.
   * The specifics of the event emission (e.g., sending data via WebSocket, triggering
   * internal application logic, etc.) are left to the implementation.
   *
   * @returns {void}
   */
  emitEvent(): void;
}
