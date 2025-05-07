import type { z } from "zod";
import type {
  IndividualGroupAPICreateSchema,
  IndividualGroupAPISchema,
  IndividualGroupCreateSchema,
  QualificationAPISchema,
  QualificationRoundAPISchema,
  QualificationRoundShrinkedAPISchema,
  QualificationSectionAPISchema,
} from "./schemas";

export type IndividualGroupAPI = z.infer<typeof IndividualGroupAPISchema>;

export type IndividualGroupAPICreate = z.infer<typeof IndividualGroupAPICreateSchema>;
export type IndividualGroupCreate = z.infer<typeof IndividualGroupCreateSchema>;

export type QualificationRoundAPI = z.infer<typeof QualificationRoundAPISchema>;
export type QualificationRoundShrinkedAPI = z.infer<typeof QualificationRoundShrinkedAPISchema>;

export type QualificationSectionAPI = z.infer<typeof QualificationSectionAPISchema>;

export type QualificationAPI = z.infer<typeof QualificationAPISchema>;
