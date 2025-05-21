import type { Competitor } from "../../../entities";
import apiClient from "../../axios/config";
import { mapToCompetitor, mapToCompetitorAPIEdit } from "./mappers";
import { CompetitorAPISchema } from "./schemas";
import type { CompetitorEdit } from "./types";

export const competitorsApi = {
  getCompetitors: async (): Promise<Competitor[]> => {
    const response = await apiClient.get("/competitors");
    const validatedResponse = CompetitorAPISchema.array().parse(response.data);
    return validatedResponse.map(mapToCompetitor);
  },
  getCompetitor: async (competitorId: number): Promise<Competitor> => {
    const response = await apiClient.get(`/competitors/${competitorId}`);
    const validatedResponse = CompetitorAPISchema.parse(response.data);
    return mapToCompetitor(validatedResponse);
  },
  putCompetitor: async (competitorId: number, data: CompetitorEdit): Promise<Competitor> => {
    const response = await apiClient.put(
      `/competitors/${competitorId}`,
      mapToCompetitorAPIEdit(data)
    );
    const validatedResponse = CompetitorAPISchema.parse(response.data);
    return mapToCompetitor(validatedResponse);
  },
  deleteCompetitor: async (competitorId: number): Promise<void> => {
    await apiClient.delete(`/competitors/${competitorId}`);
  },
};
