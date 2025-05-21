import { createContext, type Dispatch, type SetStateAction } from "react";
import type { BowClass, GroupState, Identity } from "../../../entities";

export interface GroupFilter {
  bow: BowClass | null;
  identity: Identity | null;
  state: GroupState | null;
}

export interface GroupFilterContextType {
  filter: GroupFilter;
  setFilter: Dispatch<SetStateAction<GroupFilter>>;
}

export const GroupFilterContext = createContext<GroupFilterContextType | null>(null);
