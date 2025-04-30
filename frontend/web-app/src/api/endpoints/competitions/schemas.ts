import { z } from "zod";
import { CompetitionSchema, CompetitionStage } from "../../../entities";

export const DateISOSchema = z
  .string()
  .trim()
  .regex(/^\d{4}-\d{2}-\d{2}$/)
  .date()
  .nullable();

export const CompetitionAPISchema = z.object({
  id: z.number(),
  stage: z.nativeEnum(CompetitionStage),
  start_date: DateISOSchema,
  end_date: DateISOSchema,
  is_ended: z.boolean(),
});

export const CompetitionAPICreateSchema = CompetitionAPISchema.omit({
  id: true,
  is_ended: true,
});

export const CompetitionAPIEditSchema = CompetitionAPICreateSchema.omit({
  stage: true,
});

export const CompetitionCreateSchema = CompetitionSchema.omit({
  id: true,
  isEnded: true,
});

export const CompetitionEditSchema = CompetitionCreateSchema.omit({
  stage: true,
});
