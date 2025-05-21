import { isValid, parseISO } from "date-fns";
import { z } from "zod";
import { RangeSchema, RangeType, ScoreSchema } from "../../../entities";

export const DateSchema = z
  .string()
  .trim()
  .regex(/^\d{4}-\d{2}-\d{2}$/)
  .date();

const isValidISO8601 = (value: string): boolean => {
  try {
    return isValid(parseISO(value));
  } catch {
    return false;
  }
};

export const DateTZSchema = z
  .string()
  .trim()
  .regex(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}([+-]\d{2}|Z)$/)
  .refine(isValidISO8601);

export const ShotAPISchema = z.object({
  shot_ordinal: z.number().positive(),
  score: ScoreSchema,
});

export const RangeAPISchema = z.object({
  id: z.number(),
  range_ordinal: z.number().positive(),
  is_active: z.boolean(),
  shots: ShotAPISchema.array().nullable(),
  range_score: z.number().nullable(),
});

export const RangeGroupAPISchema = z.object({
  id: z.number(),
  ranges_max_count: z.number().positive(),
  range_size: z.number().positive(),
  type: z.nativeEnum(RangeType),
  ranges: RangeAPISchema.array(),
  total_score: z.number().nullable(),
});

export const RangeAPIEditSchema = RangeAPISchema.pick({
  range_ordinal: true,
  shots: true,
});

export const RangeEditSchema = RangeSchema.pick({
  ordinal: true,
  shots: true,
});
