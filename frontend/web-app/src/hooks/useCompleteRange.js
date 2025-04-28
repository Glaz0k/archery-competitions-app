import { useMutation, useQueryClient } from "@tanstack/react-query";
import { endRange as endQualificationRange } from "../api/qualificationSections";
import { PLACE_QUERY_KEYS, SECTION_QUERY_KEYS } from "../api/queryKeys";
import { endRange as endFinalRange } from "../api/sparringPlaces";

export default function useCompleteRange(
  { placeId, sectionId, roundOrdinal },
  rangeOrdinal,
  onSuccessFn
) {
  const queryClient = useQueryClient();
  const { mutate: completeRange, isPending: isRangePending } = useMutation({
    mutationFn: async () => {
      if (placeId) {
        return await endFinalRange(placeId, rangeOrdinal);
      }
      if (sectionId && roundOrdinal) {
        return await endQualificationRange(sectionId, roundOrdinal, rangeOrdinal);
      }
      throw new Error("Invalid parameters");
    },
    onSuccess: (completedRange) => {
      if (placeId) {
        queryClient.invalidateQueries({
          queryKey: PLACE_QUERY_KEYS.element(placeId),
        });
      } else if (sectionId && roundOrdinal) {
        queryClient.invalidateQueries({
          queryKey: SECTION_QUERY_KEYS.element(sectionId),
        });
      } else {
        throw new Error("Invalid parameters");
      }
      onSuccessFn(completedRange);
    },
  });

  return { completeRange, isRangePending };
}
