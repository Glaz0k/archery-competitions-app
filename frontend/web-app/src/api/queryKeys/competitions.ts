import { CUPS_QUERY_KEYS } from "./cups";

export const COMPETITIONS_QUERY_KEYS = {
  all: ["competitions"] as const,
  allByCup: (cupId: number) => [...CUPS_QUERY_KEYS.all, cupId, ...COMPETITIONS_QUERY_KEYS.all],
  element: (competitionId: number) => [...COMPETITIONS_QUERY_KEYS.all, competitionId],
};
