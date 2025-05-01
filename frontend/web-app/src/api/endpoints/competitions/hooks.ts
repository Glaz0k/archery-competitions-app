import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { notifications } from "@mantine/notifications";
import { IndividualGroup, type Competition } from "../../../entities";
import { COMPETITIONS_QUERY_KEYS } from "../../query-keys/competitions";
import { INDIVIDUAL_GROUPS_QUERY_KEYS } from "../../query-keys/individualGroups";
import { IndividualGroupCreate } from "../individual-groups/types";
import { competitionsApi } from "./api";
import { type CompetitionEdit } from "./types";

export const useCompetition = (competitionId: number, enabled: boolean = true) => {
  return useQuery({
    queryKey: COMPETITIONS_QUERY_KEYS.element(competitionId),
    queryFn: async () => await competitionsApi.getCompetition(competitionId),
    initialData: null,
    enabled,
  });
};

export const useUpdateCompetition = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<Competition, Error, [number, CompetitionEdit]>({
    mutationFn: async ([competitionId, data]) =>
      await competitionsApi.putCompetition(competitionId, data),
    onSuccess: (editedCompetition) => {
      queryClient.invalidateQueries({
        queryKey: COMPETITIONS_QUERY_KEYS.allByCup(editedCompetition.cupId),
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

export const useEndCompetition = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<Competition, Error, number>({
    mutationFn: competitionsApi.endCompetition,
    onSuccess: (endedCompetition) => {
      queryClient.invalidateQueries({
        queryKey: COMPETITIONS_QUERY_KEYS.allByCup(endedCompetition.cupId),
      });
      queryClient.setQueryData(
        COMPETITIONS_QUERY_KEYS.element(endedCompetition.id),
        endedCompetition
      );
      onSuccess?.();
    },
  });
};

// TODO: competitors interactions

export const useCreateIndividualGroup = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<IndividualGroup, Error, [number, IndividualGroupCreate]>({
    mutationFn: async ([competitionId, data]) =>
      await competitionsApi.postIndividualGroup(competitionId, data),
    onSuccess: (createdGroup) => {
      queryClient.invalidateQueries({
        queryKey: INDIVIDUAL_GROUPS_QUERY_KEYS.allByCompetition(createdGroup.competitionId),
      });
      queryClient.setQueryData(INDIVIDUAL_GROUPS_QUERY_KEYS.element(createdGroup.id), createdGroup);
      onSuccess?.();
    },
  });
};

export const useIndividualGroups = (competitionId: number, enabled: boolean = true) => {
  return useQuery({
    queryKey: INDIVIDUAL_GROUPS_QUERY_KEYS.allByCompetition(competitionId),
    queryFn: async () => await competitionsApi.getIndidvidualGroups(competitionId),
    initialData: [],
    enabled,
  });
};
