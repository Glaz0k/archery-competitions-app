import { useMutation, useQueryClient } from "@tanstack/react-query";
import { notifications } from "@mantine/notifications";
import type { Competition } from "../../../entities";
import { COMPETITIONS_QUERY_KEYS } from "../../queryKeys/competitions";
import { competitionsApi } from "./api";
import { type CompetitionEdit } from "./types";

export const useUpdateCompetition = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<Competition, Error, [number, CompetitionEdit]>({
    mutationFn: ([competitionId, data]) => competitionsApi.putCompetition(competitionId, data),
    onSuccess: (editedCompetition) => {
      queryClient.invalidateQueries({
        queryKey: COMPETITIONS_QUERY_KEYS.all,
      });
      queryClient.setQueryData(
        COMPETITIONS_QUERY_KEYS.element(editedCompetition.id),
        editedCompetition
      );
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось изменить соревнование",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useDeleteCompetition = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<void, Error, number>({
    mutationFn: competitionsApi.deleteCompetition,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: COMPETITIONS_QUERY_KEYS.all,
      });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось удалить соревнование",
        message: error.message,
        color: "red",
      });
    },
  });
};
