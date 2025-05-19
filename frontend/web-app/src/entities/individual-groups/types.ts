import type { z } from "zod";
import type { IndividualGroupSchema } from "./schemas";

export enum Identity {
  MALES = "male",
  FEMALES = "female",
  UNITED = "united",
}

export enum GroupState {
  CREATED = "created",
  QUAL_START = "qualification_start",
  QUAL_END = "qualification_end",
  QUARTERFINAL_START = "quarterfinal_start",
  SEMIFINAL_START = "semifinal_start",
  FINAL_START = "final_start",
  COMPLETED = "completed",
}

export type IndividualGroup = z.infer<typeof IndividualGroupSchema>;
