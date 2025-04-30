import { z } from "zod";
import { CompetitionStage } from "./types";

export const CompetitionSchema = z.object({
  id: z.number(),
  stage: z.nativeEnum(CompetitionStage),
  startDate: z.date().nullable(),
  endDate: z.date().nullable(),
  isEnded: z.boolean(),
});
