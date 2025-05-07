import { z } from "zod";
import type {
  RangeAPIEditSchema,
  RangeAPISchema,
  RangeEditSchema,
  RangeGroupAPISchema,
  ShotAPISchema,
} from "./schemas";

export type ShotAPI = z.infer<typeof ShotAPISchema>;
export type RangeAPI = z.infer<typeof RangeAPISchema>;
export type RangeGroupAPI = z.infer<typeof RangeGroupAPISchema>;

export type RangeAPIEdit = z.infer<typeof RangeAPIEditSchema>;
export type RangeEdit = z.infer<typeof RangeEditSchema>;
