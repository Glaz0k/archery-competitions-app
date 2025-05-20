import { z } from "zod";
import { BowClass } from "../competitors/types";
import { GroupState, Identity } from "./types";

export const IndividualGroupSchema = z.object({
  id: z.number(),
  competitionId: z.number(),
  bow: z.nativeEnum(BowClass),
  identity: z.nativeEnum(Identity),
  state: z.nativeEnum(GroupState),
});
