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
