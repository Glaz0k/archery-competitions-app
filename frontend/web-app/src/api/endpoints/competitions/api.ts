import type { Competition } from "../../../entities";
import apiClient from "../../axios/config";
import { mapToCompetition, mapToCompetitionAPIEdit } from "./mappers";
import { CompetitionAPISchema } from "./schemas";
import type { CompetitionEdit } from "./types";

export const competitionsApi = {
  putCompetition: async (competitionId: number, data: CompetitionEdit): Promise<Competition> => {
    const response = await apiClient.put(
      `/competitions/${competitionId}`,
      mapToCompetitionAPIEdit(data)
    );
    const validatedResponse = CompetitionAPISchema.parse(response);
    return mapToCompetition(validatedResponse);
  },
  deleteCompetition: async (competitionId: number): Promise<void> => {
    await apiClient.delete(`/competitions/${competitionId}`);
  },
};
