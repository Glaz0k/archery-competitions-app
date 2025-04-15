export function createEnum(enumValues) {
  const enumObj = {};
  const valueToEntryMap = {};

  for (const [key, { value, textValue }] of Object.entries(enumValues)) {
    const entry = { value, textValue };
    enumObj[key] = entry;
    valueToEntryMap[value] = entry;
  }

  Object.defineProperty(enumObj, "valueOf", {
    value: (value) => valueToEntryMap[value],
    enumerable: false,
    configurable: false,
    writable: false,
  });

  return Object.freeze(enumObj);
}

export default createEnum;
