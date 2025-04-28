import { useQuery } from "@tanstack/react-query";
import { PLACE_QUERY_KEYS } from "../api/queryKeys";
import { getPlace } from "../api/sparringPlaces";

export default function useSparringPlace(placeId) {
  const { data: place, isFetching: isPlaceFetching } = useQuery({
    queryKey: PLACE_QUERY_KEYS.element(placeId),
    queryFn: () => getPlace(placeId),
    initialData: null,
    enabled: !!placeId,
  });
  return { place, isPlaceFetching };
}
