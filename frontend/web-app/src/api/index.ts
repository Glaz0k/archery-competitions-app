export * from "./endpoints/cups/hooks";
export * from "./endpoints/competitions/hooks";
export * from "./endpoints/individual-groups/hooks";
export * from "./endpoints/competitors/hooks";

export type { CupEdit } from "./endpoints/cups/types";
export type { CompetitionCreate, CompetitionEdit } from "./endpoints/competitions/types";
export type { IndividualGroupCreate } from "./endpoints/individual-groups/types";
export type {
  CompetitorAdd,
  CompetitorToggle,
  CompetitorEdit,
} from "./endpoints/competitors/types";

export { CupEditSchema } from "./endpoints/cups/schemas";
export { CompetitionCreateSchema, CompetitionEditSchema } from "./endpoints/competitions/schemas";
export { IndividualGroupCreateSchema } from "./endpoints/individual-groups/schemas";
export {
  CompetitorAddSchema,
  CompetitorToggleSchema,
  CompetitorEditSchema,
} from "./endpoints/competitors/schemas";

export { default as apiClient } from "./axios/config";
export { queryClient } from "./queryClient";
