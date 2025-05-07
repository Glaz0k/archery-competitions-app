import { COMPETITIONS_QUERY_KEYS } from "./competitions";

export const INDIVIDUAL_GROUPS_QUERY_KEYS = {
  all: ["individual-groups"] as const,
  allByCompetition: (competitionId: number) => [
    ...INDIVIDUAL_GROUPS_QUERY_KEYS.all,
    ...COMPETITIONS_QUERY_KEYS.element(competitionId),
  ],
  element: (groupId: number) => [...INDIVIDUAL_GROUPS_QUERY_KEYS.all, groupId],
  qualification: (groupId: number) => [
    ...INDIVIDUAL_GROUPS_QUERY_KEYS.element(groupId),
    "qualification",
  ],
};
