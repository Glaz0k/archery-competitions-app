export const DEFAULT_ENUM_VALUE = "default_enum_value";

export function defaultAndEnumValues(enumObj) {
  return [
    {
      value: DEFAULT_ENUM_VALUE,
      textValue: "Не указано",
    },
    ...Object.values(enumObj),
  ];
}
