import { CUPS_QUERY_KEYS } from "./cups";

export const COMPETITIONS_QUERY_KEYS = {
  all: ["competitions"] as const,
  allByCup: (cupId: number) => [...COMPETITIONS_QUERY_KEYS.all, ...CUPS_QUERY_KEYS.all, cupId],
  element: (competitionId: number) => [...COMPETITIONS_QUERY_KEYS.all, competitionId],
};
