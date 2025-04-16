export function isEmptyString(str) {
  return str?.trim() === "" || str == null;
}

export function isValidSeasonString(str) {
  return /^\d{4}\/\d{4}$/.test(str);
}
