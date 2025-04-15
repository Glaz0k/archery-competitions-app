export const CUP_QUERY_KEYS = {
  all: ["cups"],
  element: (id) => [...CUP_QUERY_KEYS.all, id],
};

export const COMPETITION_QUERY_KEYS = {
  all: ["competitions"],
  allByCup: (cupId) => [...CUP_QUERY_KEYS.all, cupId, ...COMPETITION_QUERY_KEYS.all],
  element: (id) => [...COMPETITION_QUERY_KEYS.all, id],
};

export const INDIVIDUAL_GROUP_QUERY_KEYS = {
  all: ["individual_groups"],
  allByCompetition: (competitionId) => [
    ...COMPETITION_QUERY_KEYS.all,
    competitionId,
    ...INDIVIDUAL_GROUP_QUERY_KEYS.all,
  ],
};
