import { z } from "zod";
import {
  BowClass,
  CompetitorCompetitionDetailSchema,
  CompetitorSchema,
  Gender,
  SportsRank,
} from "../../../entities";
import { DateSchema, DateTZSchema } from "../shared/schemas";

export const CompetitorAPISchema = z.object({
  id: z.number(),
  full_name: z.string(),
  birth_date: DateSchema,
  identity: z.nativeEnum(Gender),
  bow: z.nativeEnum(BowClass).nullable(),
  rank: z.nativeEnum(SportsRank).nullable(),
  region: z.string().nullable(),
  federation: z.string().nullable(),
  club: z.string().nullable(),
});

export const CompetitorShrinkedAPISchema = CompetitorAPISchema.pick({
  id: true,
  full_name: true,
});

export const CompetitorAPIAddSchema = z.object({
  competitor_id: z.number(),
});

export const CompetitorAddSchema = CompetitorSchema.pick({
  id: true,
});

export const CompetitorAPIEditSchema = CompetitorAPISchema.omit({
  id: true,
});

export const CompetitorEditSchema = CompetitorSchema.omit({
  id: true,
});

export const CompetitorCompetitionDetailAPISchema = z.object({
  competition_id: z.number(),
  competitor: CompetitorAPISchema,
  is_active: z.boolean(),
  created_at: DateTZSchema,
});

export const CompetitorAPIToggleSchema = CompetitorCompetitionDetailAPISchema.pick({
  is_active: true,
});

export const CompetitorToggleSchema = CompetitorCompetitionDetailSchema.pick({
  isActive: true,
});

export const CompetitorGroupDetailAPISchema = z.object({
  group_id: z.number(),
  competitor: CompetitorAPISchema,
});
