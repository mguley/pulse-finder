/**
 * Represents the result of an encryption operation.
 */
export interface IEncryptionResult {
  /**
   * The encrypted data, typically in a base64 encoded string format.
   */
  encryptedData: string;

  /**
   * The initialization vector (IV) used for the encryption, encoded as a base64 string.
   * The IV is required for the decryption process to accurately reverse the encryption.
   */
  iv: string;
}
