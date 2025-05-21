import type { z } from "zod";
import type {
  FinalAPISchema,
  FinalGridAPISchema,
  IndividualGroupAPICreateSchema,
  IndividualGroupAPISchema,
  IndividualGroupCreateSchema,
  QualificationAPISchema,
  QualificationRoundAPISchema,
  QualificationRoundShrinkedAPISchema,
  QualificationSectionAPISchema,
  QuarterfinalAPISchema,
  SemifinalAPISchema,
  ShootOutAPISchema,
  SparringAPISchema,
  SparringPlaceAPISchema,
} from "./schemas";

export type IndividualGroupAPI = z.infer<typeof IndividualGroupAPISchema>;

export type IndividualGroupAPICreate = z.infer<typeof IndividualGroupAPICreateSchema>;
export type IndividualGroupCreate = z.infer<typeof IndividualGroupCreateSchema>;

export type QualificationRoundAPI = z.infer<typeof QualificationRoundAPISchema>;
export type QualificationRoundShrinkedAPI = z.infer<typeof QualificationRoundShrinkedAPISchema>;

export type QualificationSectionAPI = z.infer<typeof QualificationSectionAPISchema>;

export type QualificationAPI = z.infer<typeof QualificationAPISchema>;

export type ShootOutAPI = z.infer<typeof ShootOutAPISchema>;

export type SparringPlaceAPI = z.infer<typeof SparringPlaceAPISchema>;

export type SparringAPI = z.infer<typeof SparringAPISchema>;

export type QuarterfinalAPI = z.infer<typeof QuarterfinalAPISchema>;
export type SemifinalAPI = z.infer<typeof SemifinalAPISchema>;
export type FinalAPI = z.infer<typeof FinalAPISchema>;
export type FinalGridAPI = z.infer<typeof FinalGridAPISchema>;
