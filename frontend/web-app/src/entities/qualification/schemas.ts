import { z } from "zod";
import { CompetitorShrinkedSchema } from "../competitors/schemas";
import { SportsRank } from "../competitors/types";
import { RangeGroupSchema } from "../score-system/schemas";

export const QualificationRoundSchema = z.object({
  sectionId: z.number(),
  ordinal: z.number().positive(),
  isActive: z.boolean(),
  rangeGroup: RangeGroupSchema,
});

export const QualificationRoundShrinkedSchema = QualificationRoundSchema.omit({
  sectionId: true,
  rangeGroup: true,
}).extend({
  totalScore: z.number().nullable(),
});

export const QualificationSectionSchema = z.object({
  id: z.number(),
  competitor: CompetitorShrinkedSchema,
  place: z.number().positive().nullable(),
  rounds: QualificationRoundShrinkedSchema.array(),
  total: z.number().nonnegative().nullable(),
  count10: z.number().nonnegative().nullable(),
  count9: z.number().nonnegative().nullable(),
  rankGained: z.nativeEnum(SportsRank).nullable(),
});

export const QualificationSchema = z.object({
  groupId: z.number(),
  distance: z.string(),
  roundCount: z.number().positive(),
  sections: QualificationSectionSchema.array(),
});
