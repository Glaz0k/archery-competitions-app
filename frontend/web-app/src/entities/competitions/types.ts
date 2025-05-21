import type { z } from "zod";
import type { CompetitionSchema } from "./schemas";

export enum CompetitionStage {
  STAGE_1 = "I",
  STAGE_2 = "II",
  STAGE_3 = "III",
  FINAL = "F",
}

export type Competition = z.infer<typeof CompetitionSchema>;
