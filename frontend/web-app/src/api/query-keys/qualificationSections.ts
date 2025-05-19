export const SECTIONS_QUERY_KEYS = {
  all: ["qualification-sections"] as const,
  element: (sectionId: number) => [...SECTIONS_QUERY_KEYS.all, sectionId],
  round: (sectionId: number, roundOrdinal: number) => [
    ...SECTIONS_QUERY_KEYS.element(sectionId),
    Number(roundOrdinal),
  ],
  rangeGroup: (sectionId: number, roundOrdinal: number) => [
    ...SECTIONS_QUERY_KEYS.round(sectionId, roundOrdinal),
    "range-group",
  ],
};
