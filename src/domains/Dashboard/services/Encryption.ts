import CryptoJS from "crypto-js";

/**
 * Represents the structure of the object returned by the encrypt method.
 * Contains the encrypted data and the initialization vector (IV) used for encryption.
 */
export interface EncryptionResult {
  encryptedData: string;
  iv: string;
}

/**
 * Defines the contract for encryption and decryption methods.
 */
export interface Decryption {
  decrypt<T>(payload: EncryptionResult): T;
}

/**
 * Represents decryption mechanism using AES-256-CBC.
 * Handles the decryption process and converts the decrypted data back to the RecentActivity type.
 */
export class AESEncryption implements Decryption {
  private readonly secretKey: string = "";

  /**
   * @param {string} secretKey - A 32-character long string used as the secret key.
   */
  constructor(secretKey: string) {
    if (secretKey.length !== 32) {
      throw new Error(`Secret key must be 32 characters long.`);
    }
    this.secretKey = secretKey;
  }

  /**
   * Decrypts the given encrypted payload using AES-256-CBC.
   *
   * @template T
   * @param {EncryptionResult} payload - The encrypted data and the IV used during encryption.
   * @returns {T} - The decrypted data, parsed as a RecentActivity object.
   * @throws Will throw an error if decryption fails.
   */
  public decrypt<T>(payload: EncryptionResult): T {
    try {
      // Convert the secret key to a WordArray for crypto-js
      const key = CryptoJS.enc.Utf8.parse(this.secretKey);

      // Convert the IV from base64 to a WordArray for crypto-js
      const iv = CryptoJS.enc.Base64.parse(payload.iv);

      // Decrypt the encrypted data using AES-256-CBC
      const decrypted = CryptoJS.AES.decrypt(payload.encryptedData, key, {
        iv: iv,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7,
      });

      // Convert the decrypted data from a WordArray to a UTF-8 string
      const decryptedText = decrypted.toString(CryptoJS.enc.Utf8);

      // Parse the decrypted string back into a RecentActivity object
      return JSON.parse(decryptedText) as T;
    } catch (e) {
      const message: string = (e as Error)?.message || "Decryption error";
      throw new Error(`Decryption failed: ${message}`);
    }
  }
}
