import { useQuery } from "@tanstack/react-query";
import { PLACE_QUERY_KEYS } from "../api/queryKeys";
import { getRangeGroup } from "../api/sparringPlaces";

export default function useSparringPlaceRangeGroup(placeId) {
  const { data: rangeGroup, isFetching: isRangeGroupFetching } = useQuery({
    queryKey: PLACE_QUERY_KEYS.rangeGroup(placeId),
    queryFn: () => getRangeGroup(placeId),
    initialData: null,
    enabled: !!placeId,
  });
  return { rangeGroup, isRangeGroupFetching };
}
