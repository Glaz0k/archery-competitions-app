export const CUP_QUERY_KEYS = {
  all: ["cups"],
  element: (cupId) => [...CUP_QUERY_KEYS.all, Number(cupId)],
};

export const COMPETITION_QUERY_KEYS = {
  all: ["competitions"],
  allByCup: (cupId) => [...COMPETITION_QUERY_KEYS.all, ...CUP_QUERY_KEYS.all, Number(cupId)],
  element: (competitionId) => [...COMPETITION_QUERY_KEYS.all, Number(competitionId)],
};

export const INDIVIDUAL_GROUP_QUERY_KEYS = {
  all: ["individual_groups"],
  allByCompetition: (competitionId) => [
    ...INDIVIDUAL_GROUP_QUERY_KEYS.all,
    ...COMPETITION_QUERY_KEYS.all,
    Number(competitionId),
  ],
  element: (groupId) => [...INDIVIDUAL_GROUP_QUERY_KEYS.all, Number(groupId)],
  qualification: (groupId) => [
    ...INDIVIDUAL_GROUP_QUERY_KEYS.all,
    Number(groupId),
    "qualification",
  ],
  finalGrid: (groupId) => [...INDIVIDUAL_GROUP_QUERY_KEYS.all, Number(groupId), "final_grid"],
};

export const COMPETITOR_QUERY_KEYS = {
  all: ["competitors"],
  allByCompetition: (competitionId) => [
    ...COMPETITOR_QUERY_KEYS.all,
    ...COMPETITION_QUERY_KEYS.all,
    Number(competitionId),
  ],
  allByGroup: (groupId) => [
    ...COMPETITOR_QUERY_KEYS.all,
    ...INDIVIDUAL_GROUP_QUERY_KEYS.all,
    Number(groupId),
  ],
};

export const SECTION_QUERY_KEYS = {
  all: ["qualification_sections"],
  element: (sectionId) => [...SECTION_QUERY_KEYS.all, Number(sectionId)],
  ranges: (sectionId, roundOrdinal) => [
    ...SECTION_QUERY_KEYS.all,
    Number(sectionId),
    Number(roundOrdinal),
  ],
};

export const PLACE_QUERY_KEYS = {
  all: ["sparring_places"],
  element: (placeId) => [...PLACE_QUERY_KEYS.all, Number(placeId)],
  rangeGroup: (placeId) => [...PLACE_QUERY_KEYS.all, Number(placeId), ...RANGE_QUERY_KEYS.all],
};

export const RANGE_QUERY_KEYS = {
  all: ["ranges"],
};
