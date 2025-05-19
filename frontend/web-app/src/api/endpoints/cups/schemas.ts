import { CupSchema } from "../../../entities";

export const CupEditSchema = CupSchema.omit({
  id: true,
});
