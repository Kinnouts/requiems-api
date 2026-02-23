import { customAlphabet } from "nanoid";

/**
 * API key generator for Requiems API
 * Generates keys in format: requiem_<24_random_chars>
 */

const ALPHABET =
	"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
const KEY_LENGTH = 24;
const nanoid = customAlphabet(ALPHABET, KEY_LENGTH);

/**
 * Generate a new API key
 * @returns Full API key (e.g., "requiem_abc123...")
 */
export function generateApiKey(): string {
	const prefix = "requiem";
	const randomPart = nanoid();
	return `${prefix}_${randomPart}`;
}

/**
 * Extract the key prefix (first 12 characters)
 * @param fullKey - Full API key
 * @returns Key prefix (e.g., "requiem_abc1")
 */
export function extractKeyPrefix(fullKey: string): string {
	return fullKey.substring(0, 12);
}

/**
 * Validate key format
 * @param key - Key to validate
 * @returns true if valid format
 */
export function isValidKeyFormat(key: string): boolean {
	// Must be requiem_ followed by 24 characters
	const pattern = /^requiem_[0-9a-zA-Z]{24}$/;
	return pattern.test(key);
}
