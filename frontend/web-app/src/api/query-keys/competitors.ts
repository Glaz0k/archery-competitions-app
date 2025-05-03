import { COMPETITIONS_QUERY_KEYS } from "./competitions";

export const COMPETITORS_QUERY_KEYS = {
  all: ["competitors"] as const,
  allByCompetition: (competitionId: number) => [
    ...COMPETITORS_QUERY_KEYS.all,
    ...COMPETITIONS_QUERY_KEYS.element(competitionId),
  ],
};
