export const PLACES_QUERY_KEYS = {
  all: ["sparring-places"] as const,
  element: (placeId: number) => [...PLACES_QUERY_KEYS.all, placeId],
  rangeGroup: (placeId: number) => [...PLACES_QUERY_KEYS.element(placeId), "range-group"],
};
