import { ShootOutSchema } from "../../../entities";
import { ShootOutAPISchema } from "../individual-groups/schemas";

export const ShootOutAPIEditSchema = ShootOutAPISchema.omit({
  id: true,
});

export const ShootOutEditSchema = ShootOutSchema.omit({
  id: true,
});
