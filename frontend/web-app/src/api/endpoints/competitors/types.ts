import type { z } from "zod";
import type {
  CompetitorAddSchema,
  CompetitorAPIAddSchema,
  CompetitorAPIEditSchema,
  CompetitorAPISchema,
  CompetitorAPIToggleSchema,
  CompetitorCompetitionDetailAPISchema,
  CompetitorEditSchema,
  CompetitorGroupDetailAPISchema,
  CompetitorShrinkedAPISchema,
  CompetitorToggleSchema,
} from "./schemas";

export type CompetitorAPI = z.infer<typeof CompetitorAPISchema>;
export type CompetitorShrinkedAPI = z.infer<typeof CompetitorShrinkedAPISchema>;

export type CompetitorAPIAdd = z.infer<typeof CompetitorAPIAddSchema>;
export type CompetitorAdd = z.infer<typeof CompetitorAddSchema>;

export type CompetitorAPIEdit = z.infer<typeof CompetitorAPIEditSchema>;
export type CompetitorEdit = z.infer<typeof CompetitorEditSchema>;

export type CompetitorCompetitionDetailAPI = z.infer<typeof CompetitorCompetitionDetailAPISchema>;

export type CompetitorToggle = z.infer<typeof CompetitorToggleSchema>;
export type CompetitorAPIToggle = z.infer<typeof CompetitorAPIToggleSchema>;

export type CompetitorGroupDetailAPI = z.infer<typeof CompetitorGroupDetailAPISchema>;
