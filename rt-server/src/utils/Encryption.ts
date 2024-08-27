import type { Cipher, Decipher } from "crypto";
import crypto from "crypto";

/**
 * EncryptionResult represents the structure of the object returned by the encrypt method,
 * containing the encrypted data and the initialization vector (IV) used for encryption.
 */
export interface EncryptionResult {
  encryptedData: string;
  iv: string;
}

/**
 * Defines the contract for encryption and decryption methods.
 */
export interface Encryption {
  /**
   * Encrypts a given plaintext string.
   *
   * @param {string} text - The plaintext string to be encrypted.
   * @returns {EncryptionResult} - Object containing the encrypted data and the IV used during encryption.
   */
  encrypt(text: string): EncryptionResult;

  /**
   * Decrypts a given encrypted string using the provided IV.
   *
   * @param {string} encryptedData - The encrypted string to be decrypted.
   * @param {string} iv - The initialization vector (IV) used during encryption.
   * @returns {string} - The decrypted plaintext string.
   */
  decrypt(encryptedData: string, iv: string): string;
}

/**
 * Provides implementation of the encryption mechanism, uses AES-256-CBC, a symmetric encryption algorithm.
 * Handles both encryption and decryption processes.
 */
export class AESEncryption implements Encryption {
  private readonly algorithm: string = "aes-256-cbc";
  private readonly secretKey: Buffer;

  /**
   * @param {string} secretKey - A 32-character long string used as the secret key.
   */
  constructor(secretKey: string) {
    this.secretKey = this.validateAndConvert(secretKey);
  }

  /**
   * Validates the secret key and converts it into a Buffer.
   *
   * @param {string} key - The secret key to validate and convert.
   * @returns {Buffer} - A Buffer containing the validated key.
   * @throws Will throw an error if the secret key is not 32 characters long.
   */
  private validateAndConvert(key: string): Buffer {
    if (key.length !== 32) {
      throw new Error(`Secret key must be 32 characters long.`);
    }
    return Buffer.from(key);
  }

  /**
   * Encrypts a given plaintext string using AES-256-CBC.
   *
   * @param {string} text - The plaintext string to be encrypted.
   * @returns {EncryptionResult} - An object containing the encrypted data and the IV used during encryption.
   * @throws Will throw an error if the encryption process fails.
   */
  public encrypt(text: string): EncryptionResult {
    try {
      const iv: Buffer = crypto.randomBytes(16);
      const cipher: Cipher = crypto.createCipheriv(
        this.algorithm,
        this.secretKey,
        iv,
      );
      let encryptedData: string = cipher.update(text, "utf8", "base64");
      encryptedData += cipher.final("base64");
      return {
        encryptedData,
        iv: iv.toString("base64"),
      };
    } catch (e) {
      const message: string =
        (e as Error)?.message || "An unknown error occurred";
      throw new Error(`Encryption failed: ${message}`);
    }
  }

  /**
   * Decrypts a given encrypted string using AES-256-CBC.
   *
   * @param {string} encryptedData - The encrypted string to be decrypted.
   * @param {string} iv - The initialization vector (IV) used during encryption, provided in base64 format.
   * @returns {string} - The decrypted plaintext string.
   * @throws Will throw an error if the decryption process fails.
   */
  public decrypt(encryptedData: string, iv: string): string {
    try {
      const decipher: Decipher = crypto.createDecipheriv(
        this.algorithm,
        this.secretKey,
        Buffer.from(iv, "base64"),
      );
      let decryptedData: string = decipher.update(
        encryptedData,
        "base64",
        "utf8",
      );
      decryptedData += decipher.final("utf8");
      return decryptedData;
    } catch (e) {
      const message: string =
        (e as Error)?.message || "An unknown error occurred";
      throw new Error(`Decryption failed: ${message}`);
    }
  }
}
