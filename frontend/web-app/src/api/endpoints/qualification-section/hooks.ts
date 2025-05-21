import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { notifications } from "@mantine/notifications";
import type { Range, RangeGroup } from "../../../entities";
import { SECTIONS_QUERY_KEYS } from "../../query-keys";
import type { RangeEdit } from "../shared/types";
import { updateRangeGroup } from "../shared/utils";
import { qualificationSectionsApi } from "./api";

export const useQualificationSection = (sectionId: number, enabled: boolean = true) => {
  return useQuery({
    queryKey: SECTIONS_QUERY_KEYS.element(sectionId),
    queryFn: async () => await qualificationSectionsApi.getSection(sectionId),
    enabled,
  });
};

export const useSectionRound = (
  sectionId: number,
  roundOrdinal: number,
  enabled: boolean = true
) => {
  return useQuery({
    queryKey: SECTIONS_QUERY_KEYS.round(sectionId, roundOrdinal),
    queryFn: async () => await qualificationSectionsApi.getRound(sectionId, roundOrdinal),
    enabled,
  });
};

export const useSectionRangeGroup = (
  sectionId: number,
  roundOrdinal: number,
  enabled: boolean = true
) => {
  return useQuery({
    queryKey: SECTIONS_QUERY_KEYS.rangeGroup(sectionId, roundOrdinal),
    queryFn: async () => await qualificationSectionsApi.getRangeGroup(sectionId, roundOrdinal),
    enabled,
  });
};

export const useEditSectionRange = (onSuccess?: () => unknown) => {
  const queryClient = useQueryClient();
  return useMutation<Range, Error, [number, number, RangeEdit]>({
    mutationFn: async ([sectionId, roundOrdinal, data]) =>
      await qualificationSectionsApi.putRange(sectionId, roundOrdinal, data),
    onSuccess: (edited, [sectionId, roundOrdinal]) => {
      queryClient.invalidateQueries({
        queryKey: SECTIONS_QUERY_KEYS.round(sectionId, roundOrdinal),
      });
      queryClient.setQueryData(
        SECTIONS_QUERY_KEYS.rangeGroup(sectionId, roundOrdinal),
        (prev: RangeGroup | undefined) => updateRangeGroup(prev, edited)
      );
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось изменить серию",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useCompleteSectionRange = (onSuccess?: () => unknown) => {
  const queryClient = useQueryClient();
  return useMutation<Range, Error, [number, number, number]>({
    mutationFn: async ([sectionId, roundOrdinal, rangeOrdinal]) =>
      await qualificationSectionsApi.endRange(sectionId, roundOrdinal, rangeOrdinal),
    onSuccess: (_, [sectionId]) => {
      queryClient.invalidateQueries({
        queryKey: SECTIONS_QUERY_KEYS.element(sectionId),
      });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось завершить серию",
        message: error.message,
        color: "red",
      });
    },
  });
};
