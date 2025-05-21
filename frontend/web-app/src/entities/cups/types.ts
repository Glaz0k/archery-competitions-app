import type { z } from "zod";
import type { CupSchema } from "./schemas";

export type Cup = z.infer<typeof CupSchema>;
