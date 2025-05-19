import type { z } from "zod";
import type {
  CompetitionAPICreateSchema,
  CompetitionAPIEditSchema,
  CompetitionAPISchema,
  CompetitionCreateSchema,
  CompetitionEditSchema,
} from "./schemas";

export type CompetitionAPI = z.infer<typeof CompetitionAPISchema>;

export type CompetitionCreate = z.infer<typeof CompetitionCreateSchema>;

export type CompetitionAPICreate = z.infer<typeof CompetitionAPICreateSchema>;

export type CompetitionEdit = z.infer<typeof CompetitionEditSchema>;

export type CompetitionAPIEdit = z.infer<typeof CompetitionAPIEditSchema>;
