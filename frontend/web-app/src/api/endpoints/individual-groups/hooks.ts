import { useMutation, useQueryClient } from "@tanstack/react-query";
import { INDIVIDUAL_GROUPS_QUERY_KEYS } from "../../query-keys/individualGroups";
import { individualGroupsApi } from "./api";

export const useDeleteIndividualGroup = (onSuccess?: () => void) => {
  const queryClient = useQueryClient();
  return useMutation<void, Error, number>({
    mutationFn: individualGroupsApi.deleteGroup,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: INDIVIDUAL_GROUPS_QUERY_KEYS.all,
      });
      onSuccess?.();
    },
  });
};
