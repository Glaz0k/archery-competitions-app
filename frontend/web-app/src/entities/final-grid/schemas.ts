import { z } from "zod";
import { CompetitorShrinkedSchema } from "../competitors/schemas";
import { RangeGroupSchema, ScoreSchema } from "../score-system/schemas";
import { SparringState } from "./types";

export const ShootOutSchema = z.object({
  id: z.number(),
  score: ScoreSchema,
  priority: z.boolean().nullable(),
});

export const SparringPlaceSchema = z.object({
  id: z.number(),
  competitor: CompetitorShrinkedSchema,
  rangeGroup: RangeGroupSchema,
  isActive: z.boolean(),
  shootOut: ShootOutSchema.nullable(),
  score: z.number(),
});

export const SparringSchema = z.object({
  id: z.number(),
  top: SparringPlaceSchema.nullable(),
  bot: SparringPlaceSchema.nullable(),
  state: z.nativeEnum(SparringState),
});

export const FinalSchema = z.object({
  sparringGold: SparringSchema.nullable(),
  sparringBronze: SparringSchema.nullable(),
});

export const SemifinalSchema = z.object({
  sparring5: SparringSchema.nullable(),
  sparring6: SparringSchema.nullable(),
});

export const QuarterfinalSchema = z.object({
  sparring1: SparringSchema.nullable(),
  sparring2: SparringSchema.nullable(),
  sparring3: SparringSchema.nullable(),
  sparring4: SparringSchema.nullable(),
});

export const FinalGridSchema = z.object({
  groupId: z.number(),
  quarterfinal: QuarterfinalSchema,
  semifinal: SemifinalSchema.nullable(),
  final: FinalSchema.nullable(),
});
