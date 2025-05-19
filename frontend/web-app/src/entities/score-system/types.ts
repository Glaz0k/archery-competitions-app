import type { z } from "zod";
import type { RangeGroupSchema, RangeSchema, ScoreSchema, ShotSchema } from "./schemas";

export enum RangeType {
  ONE_TEN = "1-10",
  SIX_TEN = "6-10",
}

export type Score = z.infer<typeof ScoreSchema>;
export type Shot = z.infer<typeof ShotSchema>;
export type Range = z.infer<typeof RangeSchema>;
export type RangeGroup = z.infer<typeof RangeGroupSchema>;
