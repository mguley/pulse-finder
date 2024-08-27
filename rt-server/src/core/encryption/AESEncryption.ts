import type { Cipher, Decipher } from "crypto";
import crypto from "crypto";
import type { IEncryption } from "../interfaces/encryption/IEncryption";
import {
  DEFAULT_ENCRYPTION_ALGO,
  DEFAULT_SECRET_KEY,
} from "../interfaces/config/IConfig";
import type { IEncryptionResult } from "../interfaces/encryption/IEncryptionResult";

/**
 * Provides encryption and decryption capabilities using the AES-256-CBC algorithm.
 */
export class AESEncryption implements IEncryption {
  private readonly algorithm: string;
  private readonly secretKey: Buffer;

  /**
   * @param {string} [algorithm] - The encryption algorithm to use.
   * @param {string} [secretKey] - The secret key for encryption.
   */
  constructor(algorithm?: string, secretKey?: string) {
    this.algorithm = algorithm ?? DEFAULT_ENCRYPTION_ALGO;
    this.secretKey = secretKey ? this.validate(secretKey) : DEFAULT_SECRET_KEY;
  }

  /**
   * Validates the provided secret key.
   * The secret key must be 32 characters long to be used for AES-256 encryption.
   *
   * @param {string} key - The secret key to validate.
   * @returns {Buffer} - The validated secret key as a Buffer.
   * @throws Will throw an error if the secret key is not 32 characters long.
   */
  private validate(key: string): Buffer {
    if (key.length !== 32) {
      throw new Error(`Secret key must be 32 characters long.`);
    }
    return Buffer.from(key);
  }

  /**
   * Encrypts a given plaintext string using AES-256-CBC.
   *
   * @param {string} text - The plaintext string to encrypt.
   * @returns {IEncryptionResult} - The result of the encryption, including the encrypted data and IV.
   * @throws Will throw an error if the encryption process fails.
   */
  public encrypt(text: string): IEncryptionResult {
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
   * @param {string} encryptedData - The encrypted string to decrypt.
   * @param {string} iv - The initialization vector (IV) used during encryption, provided in base64 format.
   * @returns {string} The decrypted plaintext string.
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
