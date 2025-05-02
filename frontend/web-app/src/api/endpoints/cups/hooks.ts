import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { notifications } from "@mantine/notifications";
import type { Competition, Cup } from "../../../entities";
import { COMPETITIONS_QUERY_KEYS } from "../../query-keys/competitions";
import { CUPS_QUERY_KEYS } from "../../query-keys/cups";
import { cupsApi } from "./api";
import type { CupEdit } from "./types";

export const useCups = (enabled: boolean = true) => {
  return useQuery({
    queryKey: CUPS_QUERY_KEYS.all,
    queryFn: cupsApi.getCups,
    initialData: [],
    enabled,
  });
};

export const useCreateCup = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<Cup, Error, CupEdit>({
    mutationFn: cupsApi.postCup,
    onSuccess: (createdCup) => {
      queryClient.setQueryData(CUPS_QUERY_KEYS.all, (old: Cup[]) => {
        return [createdCup, ...old];
      });
      queryClient.setQueryData(CUPS_QUERY_KEYS.element(createdCup.id), createdCup);
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось создать кубок",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useDeleteCup = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<void, Error, number>({
    mutationFn: cupsApi.deleteCup,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: CUPS_QUERY_KEYS.all,
      });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось удалить кубок",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useCup = (cupId: number, enabled: boolean = true) => {
  return useQuery({
    queryKey: CUPS_QUERY_KEYS.element(cupId),
    queryFn: async () => await cupsApi.getCup(cupId),
    initialData: null,
    enabled,
  });
};

export const useUpdateCup = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<Cup, Error, Parameters<typeof cupsApi.putCup>>({
    mutationFn: async ([cupId, data]) => await cupsApi.putCup(cupId, data),
    onSuccess: (editedCup) => {
      queryClient.setQueryData(CUPS_QUERY_KEYS.element(editedCup.id), { ...editedCup });
      queryClient.invalidateQueries({ queryKey: CUPS_QUERY_KEYS.element(editedCup.id) });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось изменить кубок",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useCreateCompetition = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<Competition, Error, Parameters<typeof cupsApi.postCompetiton>>({
    mutationFn: async ([cupId, data]) => await cupsApi.postCompetiton(cupId, data),
    onSuccess: (createdCompetition, [cupId]) => {
      queryClient.setQueryData(COMPETITIONS_QUERY_KEYS.allByCup(cupId), (old: Competition[]) => {
        return [createdCompetition, ...old];
      });
      onSuccess?.();
    },
    onError: (error) => {
      notifications.show({
        title: "Не удалось создать соревнование",
        message: error.message,
        color: "red",
      });
    },
  });
};

export const useCompetitions = (cupId: number, enabled: boolean = true) => {
  return useQuery({
    queryKey: COMPETITIONS_QUERY_KEYS.allByCup(cupId),
    queryFn: async () => await cupsApi.getCompetitons(cupId),
    initialData: [],
    enabled,
  });
};
