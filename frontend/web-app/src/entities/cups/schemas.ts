import { z } from "zod";
import { SEASON_REGEX } from "./constants";

export const CupSchema = z.object({
  id: z.number(),
  title: z.string().trim().nonempty(),
  address: z.string().trim().nullable(),
  season: z.string().trim().regex(SEASON_REGEX).nullable(),
});
