import type { z } from "zod";
import type {
  CompetitorAddSchema,
  CompetitorAPIAddSchema,
  CompetitorAPISchema,
  CompetitorAPIToggleSchema,
  CompetitorCompetitionDetailAPISchema,
  CompetitorToggleSchema,
} from "./schemas";

export type CompetitorAPI = z.infer<typeof CompetitorAPISchema>;

export type CompetitorAPIAdd = z.infer<typeof CompetitorAPIAddSchema>;

export type CompetitorAdd = z.infer<typeof CompetitorAddSchema>;

export type CompetitorCompetitionDetailAPI = z.infer<typeof CompetitorCompetitionDetailAPISchema>;

export type CompetitorToggle = z.infer<typeof CompetitorToggleSchema>;

export type CompetitorAPIToggle = z.infer<typeof CompetitorAPIToggleSchema>;
