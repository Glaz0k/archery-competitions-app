export function createEnum(enumValues) {
  const enumObj = {};
  const valueToKeyMap = {};

  for (const [key, { value, textValue }] of Object.entries(enumValues)) {
    enumObj[key] = { value, textValue };
    valueToKeyMap[value] = key;
  }

  Object.defineProperty(enumObj, "valueOf", {
    value: (value) => enumObj[valueToKeyMap[value]],
    enumerable: false,
    configurable: false,
    writable: false,
  });

  return Object.freeze(enumObj);
}

export default createEnum;
