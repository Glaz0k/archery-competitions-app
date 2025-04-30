export const CUPS_QUERY_KEYS = {
  all: ["cups"] as const,
  element: (cupId: number) => [...CUPS_QUERY_KEYS.all, cupId] as const,
};
