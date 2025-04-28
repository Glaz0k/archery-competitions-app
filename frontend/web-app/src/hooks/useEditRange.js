import { useMutation, useQueryClient } from "@tanstack/react-query";
import { putRange as putQualificationRange } from "../api/qualificationSections";
import { PLACE_QUERY_KEYS, SECTION_QUERY_KEYS } from "../api/queryKeys";
import { putRange as putFinalRange } from "../api/sparringPlaces";

export default function useEditRange(
  { placeId, sectionId, roundOrdinal },
  rangeOrdinal,
  onSuccessFn
) {
  const queryClient = useQueryClient();
  const { mutateAsync: asyncEditRange, isPending: isRangePending } = useMutation({
    mutationFn: async (changedShots) => {
      if (placeId) {
        return await putFinalRange(placeId, { rangeOrdinal: rangeOrdinal, shots: changedShots });
      }
      if (sectionId && roundOrdinal) {
        return await putQualificationRange(sectionId, roundOrdinal, {
          rangeOrdinal: rangeOrdinal,
          shots: changedShots,
        });
      }
      throw new Error("Invalid parameters");
    },
    onSuccess: (editedRange) => {
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
      onSuccessFn(editedRange);
    },
  });

  return { asyncEditRange, isRangePending };
}
