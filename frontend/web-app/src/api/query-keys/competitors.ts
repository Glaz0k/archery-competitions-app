import { COMPETITIONS_QUERY_KEYS } from "./competitions";
import { INDIVIDUAL_GROUPS_QUERY_KEYS } from "./individualGroups";

export const COMPETITORS_QUERY_KEYS = {
  all: ["competitors"] as const,
  allByCompetition: (competitionId: number) => [
    ...COMPETITORS_QUERY_KEYS.all,
    ...COMPETITIONS_QUERY_KEYS.element(competitionId),
  ],
  allByGroup: (groupId: number) => [
    ...COMPETITORS_QUERY_KEYS.all,
    ...INDIVIDUAL_GROUPS_QUERY_KEYS.element(groupId),
  ],
  element: (competitorId: number) => [...COMPETITORS_QUERY_KEYS.all, competitorId],
};
