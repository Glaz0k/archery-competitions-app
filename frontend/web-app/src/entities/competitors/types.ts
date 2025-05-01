import type { z } from "zod";
import type {
  CompetitorCompetitionDetailSchema,
  CompetitorSchema,
  CompetitorShrinkedSchema,
} from "./schemas";

export enum Gender {
  MALE = "male",
  FEMALE = "female",
}

export enum BowClass {
  CLASSIC = "classic",
  BLOCK = "block",
  CLASSIC_NEWBIE = "classic_newbie",
  CLASSIC_3D = "3D_classic",
  COMPOUND_3D = "3D_compound",
  LONG_3D = "3D_long",
  PERIPHERAL = "peripheral",
  PERIPHERAL_WITH_RING = "peripheral_with_ring",
}

export enum SportsRank {
  MASTER_MERITED = "merited_master",
  MASTER_INTERNATIONAL = "master_international",
  MASTER = "master",
  MASTER_CANDIDATE = "candidate_for_master",
  CLASS_1 = "first_class",
  CLASS_2 = "second_class",
  CLASS_3 = "third_class",
}

export type Competitor = z.infer<typeof CompetitorSchema>;

export type CompetitorShrinked = z.infer<typeof CompetitorShrinkedSchema>;

export type CompetitorCompetitionDetail = z.infer<typeof CompetitorCompetitionDetailSchema>;
