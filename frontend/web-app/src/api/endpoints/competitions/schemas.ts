import { z } from "zod";
import { CompetitionSchema, CompetitionStage } from "../../../entities";
import { DateSchema } from "../shared/schemas";

export const CompetitionAPISchema = z.object({
  id: z.number(),
  cup_id: z.number(),
  stage: z.nativeEnum(CompetitionStage),
  start_date: DateSchema.nullable(),
  end_date: DateSchema.nullable(),
  is_ended: z.boolean(),
});

export const CompetitionAPICreateSchema = CompetitionAPISchema.omit({
  id: true,
  cup_id: true,
  is_ended: true,
});

export const CompetitionCreateSchema = CompetitionSchema.omit({
  id: true,
  cupId: true,
  isEnded: true,
});

export const CompetitionAPIEditSchema = CompetitionAPICreateSchema.omit({
  stage: true,
});

export const CompetitionEditSchema = CompetitionCreateSchema.omit({
  stage: true,
});
