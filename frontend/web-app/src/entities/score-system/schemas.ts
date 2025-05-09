import { z } from "zod";
import { RangeType } from "./types";

export const ScoreSchema = z
  .string()
  .trim()
  .regex(/^(M|[1-9]|10|X)$/)
  .nullable();

export const ShotSchema = z.object({
  ordinal: z.number().positive(),
  score: ScoreSchema,
});

export const RangeSchema = z.object({
  id: z.number(),
  ordinal: z.number().positive(),
  isActive: z.boolean(),
  shots: ShotSchema.array().nullable(),
  score: z.number().nullable(),
});

export const RangeGroupSchema = z.object({
  id: z.number(),
  rangesMaxCount: z.number().positive(),
  rangeSize: z.number().positive(),
  type: z.nativeEnum(RangeType),
  ranges: RangeSchema.array(),
  totalScore: z.number().nullable(),
});
