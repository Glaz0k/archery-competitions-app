export const CUP_QUERY_KEYS = {
  all: ["cups"],
  element: (cupId) => [...CUP_QUERY_KEYS.all, Number(cupId)],
};

export const COMPETITION_QUERY_KEYS = {
  all: ["competitions"],
  allByCup: (cupId) => [...CUP_QUERY_KEYS.all, Number(cupId), ...COMPETITION_QUERY_KEYS.all],
  element: (competitionId) => [...COMPETITION_QUERY_KEYS.all, Number(competitionId)],
};

export const INDIVIDUAL_GROUP_QUERY_KEYS = {
  all: ["individual_groups"],
  allByCompetition: (competitionId) => [
    ...COMPETITION_QUERY_KEYS.all,
    Number(competitionId),
    ...INDIVIDUAL_GROUP_QUERY_KEYS.all,
  ],
};

export const COMPETITOR_QUERY_KEYS = {
  all: ["competitors"],
  allByCompetition: (competitionId) => [...COMPETITOR_QUERY_KEYS.all, Number(competitionId)],
};
