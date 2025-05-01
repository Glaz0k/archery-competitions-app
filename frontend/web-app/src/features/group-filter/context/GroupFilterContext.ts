import { createContext } from "react";
import type { BowClass, GroupState, Identity } from "../../../entities";

export interface GroupFilter {
  bow: BowClass | undefined;
  identity: Identity | undefined;
  state: GroupState | undefined;
}

export interface GroupFilterContextType {
  filter: GroupFilter;
  setFilter: (filter: GroupFilter) => void;
}

export const GroupFilterContext = createContext<GroupFilterContextType | null>(null);
