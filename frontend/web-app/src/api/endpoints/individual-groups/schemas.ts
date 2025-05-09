import { z } from "zod";
import {
  BowClass,
  Gender,
  GroupState,
  IndividualGroupSchema,
  ScoreSchema,
  SparringState,
  SportsRank,
} from "../../../entities";
import { CompetitorShrinkedAPISchema } from "../competitors/schemas";
import { RangeGroupAPISchema } from "../shared/schemas";

export const IndividualGroupAPISchema = z.object({
  id: z.number(),
  competition_id: z.number(),
  bow: z.nativeEnum(BowClass),
  identity: z.nativeEnum(Gender).nullable(),
  state: z.nativeEnum(GroupState),
});

export const IndividualGroupAPICreateSchema = IndividualGroupAPISchema.pick({
  bow: true,
  identity: true,
});

export const IndividualGroupCreateSchema = IndividualGroupSchema.pick({
  bow: true,
  identity: true,
});

export const QualificationRoundAPISchema = z.object({
  section_id: z.number(),
  round_ordinal: z.number().positive(),
  is_active: z.boolean(),
  range_group: RangeGroupAPISchema,
});

export const QualificationRoundShrinkedAPISchema = QualificationRoundAPISchema.omit({
  section_id: true,
  range_group: true,
}).extend({
  total_score: z.number().nullable(),
});

export const QualificationSectionAPISchema = z.object({
  id: z.number(),
  competitor: CompetitorShrinkedAPISchema,
  place: z.number().positive().nullable(),
  rounds: QualificationRoundShrinkedAPISchema.array(),
  total: z.number().nonnegative().nullable(),
  "10_s": z.number().nonnegative().nullable(),
  "9_s": z.number().nonnegative().nullable(),
  rank_gained: z.nativeEnum(SportsRank).nullable(),
});

export const QualificationAPISchema = z.object({
  group_id: z.number(),
  distance: z.string(),
  round_count: z.number().positive(),
  sections: QualificationSectionAPISchema.array(),
});

export const ShootOutAPISchema = z.object({
  id: z.number(),
  score: ScoreSchema,
  priority: z.boolean().nullable(),
});

export const SparringPlaceAPISchema = z.object({
  id: z.number(),
  competitor: CompetitorShrinkedAPISchema,
  range_group: RangeGroupAPISchema,
  is_active: z.boolean(),
  shoot_out: ShootOutAPISchema.nullable(),
  sparring_score: z.number().nonnegative(),
});

export const SparringAPISchema = z.object({
  id: z.number(),
  top_place: SparringPlaceAPISchema.nullable(),
  bot_place: SparringPlaceAPISchema.nullable(),
  state: z.nativeEnum(SparringState),
});

export const QuarterfinalAPISchema = z.object({
  sparring_1: SparringAPISchema.nullable(),
  sparring_2: SparringAPISchema.nullable(),
  sparring_3: SparringAPISchema.nullable(),
  sparring_4: SparringAPISchema.nullable(),
});

export const SemifinalAPISchema = z.object({
  sparring_5: SparringAPISchema.nullable(),
  sparring_6: SparringAPISchema.nullable(),
});

export const FinalAPISchema = z.object({
  sparring_gold: SparringAPISchema.nullable(),
  sparring_bronze: SparringAPISchema.nullable(),
});

export const FinalGridAPISchema = z.object({
  group_id: z.number(),
  quarterfinal: QuarterfinalAPISchema,
  semifinal: SemifinalAPISchema.nullable(),
  final: FinalAPISchema.nullable(),
});
