import { z } from "zod";
import { BowClass, Gender, GroupState, IndividualGroupSchema, SportsRank } from "../../../entities";
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
