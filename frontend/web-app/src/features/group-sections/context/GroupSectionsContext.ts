import { createContext, type Dispatch, type SetStateAction } from "react";
import type { Competition, Cup, IndividualGroup } from "../../../entities";

export interface GroupSections {
  info: {
    cup: Cup | null;
    competition: Competition | null;
    group: IndividualGroup | null;
  };
  controls: {
    onRefresh: (() => unknown) | undefined;
    onComplete: (() => unknown) | undefined;
    onExport: (() => unknown) | undefined;
  };
}

export interface GroupSectionsContextType {
  context: GroupSections;
  setContext: Dispatch<SetStateAction<GroupSections>>;
}

export const GroupSectionsContext = createContext<GroupSectionsContextType | null>(null);
