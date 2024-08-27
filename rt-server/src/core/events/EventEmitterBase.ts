import type { IEventEmitter } from "../interfaces/event/IEventEmitter";

/**
 * Serves as an abstract base class for all event emitters.
 *
 * @abstract
 * @implements {IEventEmitter}
 */
export abstract class EventEmitterBase implements IEventEmitter {
  /**
   * Responsible for emitting events.
   *
   * @abstract
   * @returns {void}
   */
  abstract emitEvent(): void;
}
