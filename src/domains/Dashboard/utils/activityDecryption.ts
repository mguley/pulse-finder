import type { EncryptionResult } from "../services/Encryption";
import { AESEncryption, type Decryption } from "../services/Encryption";

const secretKey: string = "12345678901234567890123456789012";
const decrypt: Decryption = new AESEncryption(secretKey);

/**
 * Type guard to check if the incoming data is an EncryptionResult.
 *
 * @param {object} data - The data to check.
 * @returns {boolean} - True if the data is an EncryptionResult, false otherwise.
 */
export function isEncryptionResult(data: object): data is EncryptionResult {
  return (
    data &&
    typeof data === "object" &&
    typeof (data as EncryptionResult)?.encryptedData === "string" &&
    typeof (data as EncryptionResult)?.iv === "string"
  );
}

/**
 * Handles decryption if the data is encrypted.
 *
 * @template T
 * @param {T | EncryptionResult} data - The data to decrypt if needed.
 * @returns {T} - The decrypted or plain data.
 */
export function handleDecryption<T>(data: T | EncryptionResult): T {
  return isEncryptionResult(data) ? (decrypt.decrypt(data) as T) : data;
}
