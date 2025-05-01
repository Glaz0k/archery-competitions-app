import type { z } from "zod";
import type {
  IndividualGroupAPICreateSchema,
  IndividualGroupAPISchema,
  IndividualGroupCreateSchema,
} from "./schemas";

export type IndividualGroupAPI = z.infer<typeof IndividualGroupAPISchema>;

export type IndividualGroupAPICreate = z.infer<typeof IndividualGroupAPICreateSchema>;

export type IndividualGroupCreate = z.infer<typeof IndividualGroupCreateSchema>;
