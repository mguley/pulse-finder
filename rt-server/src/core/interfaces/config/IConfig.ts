/**
 * Defines the time interval (in milliseconds) for emitting events.
 *
 * @constant
 * @type {number}
 */
export const TIME_INTERVAL: number = 5000;

/**
 * The default encryption algorithm, which is `aes-256-cbc`.
 *
 * @constant
 * @type {string}
 */
export const DEFAULT_ENCRYPTION_ALGO: string = "aes-256-cbc";

/**
 * Represents the default secret key used for encryption. It must be 32 bytes long.
 *
 * @constant
 * @type {Buffer}
 */
export const DEFAULT_SECRET_KEY: Buffer = Buffer.from(
  "12345678901234567890123456789012",
);
