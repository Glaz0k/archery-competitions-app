import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { notifications } from "@mantine/notifications";
import type { Competition, CompetitorCompetitionDetail, IndividualGroup } from "../../../entities";
import { COMPETITIONS_QUERY_KEYS } from "../../query-keys/competitions";
import { COMPETITORS_QUERY_KEYS } from "../../query-keys/competitors";
import { INDIVIDUAL_GROUPS_QUERY_KEYS } from "../../query-keys/individualGroups";
import type { CompetitorAdd, CompetitorToggle } from "../competitors/types";
import type { IndividualGroupCreate } from "../individual-groups/types";
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
    onError: (error) => {
      notifications.show({
        title: "Не удалось завершить соревнование",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useAddCompetitorToCompetition = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<CompetitorCompetitionDetail, Error, [number, CompetitorAdd]>({
    mutationFn: async ([competitionId, data]) =>
      await competitionsApi.postCompetitor(competitionId, data),
    onSuccess: (_, [competitionId]) => {
      queryClient.invalidateQueries({
        queryKey: COMPETITORS_QUERY_KEYS.allByCompetition(competitionId),
      });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось добавить участника в соревнование",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useCompetitionCompetitors = (competitionId: number, enabled: boolean = true) => {
  return useQuery({
    queryKey: COMPETITORS_QUERY_KEYS.allByCompetition(competitionId),
    queryFn: async () => await competitionsApi.getCompetitors(competitionId),
    initialData: [],
    enabled,
  });
};

export const useToggleCompetitor = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<CompetitorCompetitionDetail, Error, [number, number, CompetitorToggle]>({
    mutationFn: async ([competitionId, competitorId, data]) =>
      await competitionsApi.putCompetitor(competitionId, competitorId, data),
    onSuccess: (toggled) => {
      queryClient.setQueryData(
        COMPETITORS_QUERY_KEYS.allByCompetition(toggled.competitionId),
        (old: CompetitorCompetitionDetail[]) =>
          old.map((detail) => (detail.competitor.id === toggled.competitor.id ? toggled : detail))
      );
      onSuccess?.();
    },
    onError: (error, [, , { isActive }]) => {
      notifications.show({
        title: `Не удалось ${isActive ? "активировать" : "деактивировать"} участника`,
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useRemoveCompetitorFromCompetition = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<void, Error, [number, number]>({
    mutationFn: async ([competitionId, competitorId]) =>
      await competitionsApi.deleteCompetitor(competitionId, competitorId),
    onSuccess: (_, [competitionId, competitorId]) => {
      queryClient.setQueryData(
        COMPETITORS_QUERY_KEYS.allByCompetition(competitionId),
        (old: CompetitorCompetitionDetail[]) =>
          old.filter((detail) => detail.competitor.id !== competitorId)
      );
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: `Не удалось исключить участника`,
        message: error.message,
        color: "red",
      });
    },
  });
};

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
    onError: (error) => {
      notifications.show({
        title: "Не удалось создать дивизион",
        message: error.message,
        color: "red",
      });
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
