import type { EncryptionResult } from "../services/Encryption";
import { AESEncryption } from "../services/Encryption";
import type { RecentActivity } from "../services/recentActivity/types";

const secretKey = "12345678901234567890123456789012";
const decrypt = new AESEncryption(secretKey);

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
 * Handles decryption if the activity is encrypted.
 *
 * @param {RecentActivity | EncryptionResult} activity - The activity data to decrypt if needed.
 * @returns {RecentActivity} - The decrypted or plain activity data.
 */
export function handleActivityDecryption(
  activity: RecentActivity | EncryptionResult,
): RecentActivity {
  return isEncryptionResult(activity) ? decrypt.decrypt(activity) : activity;
}
