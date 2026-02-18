import { customAlphabet } from "nanoid";

/**
 * API key generator for Requiems API
 * Generates keys in format: requiem_<24_random_chars>
 */
export class ApiKeyGenerator {
	private static readonly ALPHABET =
		"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
	private static readonly KEY_LENGTH = 24;
	private static readonly nanoid = customAlphabet(
		ApiKeyGenerator.ALPHABET,
		ApiKeyGenerator.KEY_LENGTH,
	);

	/**
	 * Generate a new API key
	 * @returns Full API key (e.g., "requiem_abc123...")
	 */
	static generate(): string {
		const prefix = "requiem";
		const randomPart = this.nanoid();
		return `${prefix}_${randomPart}`;
	}

	/**
	 * Extract the key prefix (first 12 characters)
	 * @param fullKey - Full API key
	 * @returns Key prefix (e.g., "requiem_abc1")
	 */
	static extractPrefix(fullKey: string): string {
		return fullKey.substring(0, 12);
	}

	/**
	 * Validate key format
	 * @param key - Key to validate
	 * @returns true if valid format
	 */
	static isValidFormat(key: string): boolean {
		// Must be requiem_ followed by 24 characters
		const pattern = /^requiem_[0-9a-zA-Z]{24}$/;
		return pattern.test(key);
	}
}
