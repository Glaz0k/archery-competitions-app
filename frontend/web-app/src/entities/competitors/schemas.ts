import { z } from "zod";
import { BowClass, Gender, SportsRank } from "./types";

export const CompetitorSchema = z.object({
  id: z.number(),
  fullName: z.string().nonempty(),
  birthDate: z.date(),
  identity: z.nativeEnum(Gender),
  bow: z.nativeEnum(BowClass).nullable(),
  rank: z.nativeEnum(SportsRank).nullable(),
  region: z.string().nullable(),
  federation: z.string().nullable(),
  club: z.string().nullable(),
});

export const CompetitorShrinkedSchema = CompetitorSchema.pick({
  id: true,
  fullName: true,
});

export const CompetitorCompetitionDetailSchema = z.object({
  competitionId: z.number(),
  competitor: CompetitorSchema,
  isActive: z.boolean(),
  createdAt: z.date(),
});

export const CompetitorGroupDetailSchema = z.object({
  groupId: z.number(),
  competitor: CompetitorSchema,
});
