import { z } from "zod";
import { BowClass, Gender, GroupState, IndividualGroupSchema } from "../../../entities";

export const IndividualGroupAPISchema = z.object({
  id: z.number(),
  competition_id: z.number(),
  bow: z.nativeEnum(BowClass),
  identity: z.nativeEnum(Gender).nullable(),
  state: z.nativeEnum(GroupState),
});

export const IndividualGroupAPICreateSchema = IndividualGroupAPISchema.pick({
  bow: true,
  identity: true,
});

export const IndividualGroupCreateSchema = IndividualGroupSchema.pick({
  bow: true,
  identity: true,
});
