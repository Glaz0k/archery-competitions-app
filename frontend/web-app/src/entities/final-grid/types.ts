import type { z } from "zod";
import type {
  FinalGridSchema,
  FinalSchema,
  QuarterfinalSchema,
  SemifinalSchema,
  ShootOutSchema,
  SparringPlaceSchema,
  SparringSchema,
} from "./schemas";

export enum SparringState {
  ONGOING = "ongoing",
  TOP_WIN = "top_win",
  BOT_WIN = "bot_win",
}

export type ShootOut = z.infer<typeof ShootOutSchema>;

export type SparringPlace = z.infer<typeof SparringPlaceSchema>;
export type Sparring = z.infer<typeof SparringSchema>;

export type Final = z.infer<typeof FinalSchema>;
export type Semifinal = z.infer<typeof SemifinalSchema>;
export type Quarterfinal = z.infer<typeof QuarterfinalSchema>;
export type FinalGrid = z.infer<typeof FinalGridSchema>;
