import type {
  QualificationRound,
  QualificationSection,
  Range,
  RangeGroup,
} from "../../../entities";
import apiClient from "../../axios/config";
import { mapToQualificationRound, mapToSection } from "../individual-groups/mappers";
import {
  QualificationRoundAPISchema,
  QualificationSectionAPISchema,
} from "../individual-groups/schemas";
import { mapToRange, mapToRangeAPIEdit, mapToRangeGroup } from "../shared/mappers";
import { RangeAPISchema, RangeGroupAPISchema } from "../shared/schemas";
import type { RangeEdit } from "../shared/types";

export const qualificationSectionsApi = {
  getSection: async (sectionId: number): Promise<QualificationSection> => {
    const response = await apiClient.get(`/qualification_sections/${sectionId}`);
    const validatedResponse = QualificationSectionAPISchema.parse(response.data);
    return mapToSection(validatedResponse);
  },
  getRound: async (sectionId: number, roundOrdinal: number): Promise<QualificationRound> => {
    const response = await apiClient.get(
      `/qualification_sections/${sectionId}/rounds/${roundOrdinal}`
    );
    const validatedResponse = QualificationRoundAPISchema.parse(response.data);
    return mapToQualificationRound(validatedResponse);
  },
  getRangeGroup: async (sectionId: number, roundOrdinal: number): Promise<RangeGroup> => {
    const response = await apiClient.get(
      `/qualification_sections/${sectionId}/rounds/${roundOrdinal}/ranges`
    );
    const validatedResponse = RangeGroupAPISchema.parse(response.data);
    return mapToRangeGroup(validatedResponse);
  },
  putRange: async (sectionId: number, roundOrdinal: number, data: RangeEdit): Promise<Range> => {
    const response = await apiClient.put(
      `/qualification_sections/${sectionId}/rounds/${roundOrdinal}/ranges`,
      mapToRangeAPIEdit(data)
    );
    const validatedResponse = RangeAPISchema.parse(response.data);
    return mapToRange(validatedResponse);
  },
  endRange: async (
    sectionId: number,
    roundOrdinal: number,
    rangeOrdinal: number
  ): Promise<Range> => {
    const response = await apiClient.post(
      `/qualification_sections/${sectionId}/rounds/${roundOrdinal}/ranges/${rangeOrdinal}/end`
    );
    const validatedResponse = RangeAPISchema.parse(response.data);
    return mapToRange(validatedResponse);
  },
};
