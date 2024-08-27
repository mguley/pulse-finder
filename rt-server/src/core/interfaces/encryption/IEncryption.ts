import type { IEncryptionResult } from "./IEncryptionResult";

/**
 * Defines the contract for encryption/decryption mechanism.
 */
export interface IEncryption {
  /**
   * Encrypts a given plaintext string.
   *
   * @param {string} text - The plaintext string to be encrypted.
   * @returns {IEncryptionResult} - Object containing the encrypted data and the IV used during encryption.
   */
  encrypt(text: string): IEncryptionResult;

  /**
   * Decrypts a given encrypted string using the provided IV.
   *
   * @param {string} encryptedData - The encrypted string to be decrypted.
   * @param {string} iv - The initialization vector (IV) used during encryption.
   * @returns {string} - The decrypted plaintext string.
   */
  decrypt(encryptedData: string, iv: string): string;
}
