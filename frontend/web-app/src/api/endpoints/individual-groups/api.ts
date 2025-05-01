import apiClient from "../../axios/config";

export const individualGroupsApi = {
  // getGroup
  deleteGroup: async (groupId: number): Promise<void> => {
    await apiClient.delete(`/individual_groups/${groupId}`);
  },
};
