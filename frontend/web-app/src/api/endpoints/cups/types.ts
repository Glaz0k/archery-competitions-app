import type { z } from "zod";
import type { CupEditSchema } from "./schemas";

export type CupEdit = z.infer<typeof CupEditSchema>;
