import type { z } from "zod";
import type {
  QualificationRoundSchema,
  QualificationRoundShrinkedSchema,
  QualificationSchema,
  QualificationSectionSchema,
} from "./schemas";

export type QualificationRound = z.infer<typeof QualificationRoundSchema>;
export type QualificationRoundShrinked = z.infer<typeof QualificationRoundShrinkedSchema>;
export type QualificationSection = z.infer<typeof QualificationSectionSchema>;
export type Qualification = z.infer<typeof QualificationSchema>;
