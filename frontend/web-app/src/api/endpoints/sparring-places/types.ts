import type { z } from "zod";
import type { ShootOutAPIEditSchema, ShootOutEditSchema } from "./schemas";

export type ShootOutAPIEdit = z.infer<typeof ShootOutAPIEditSchema>;
export type ShootOutEdit = z.infer<typeof ShootOutEditSchema>;
