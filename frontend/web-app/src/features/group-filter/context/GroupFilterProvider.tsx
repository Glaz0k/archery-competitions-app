import { useState, type ReactNode } from "react";
import { GroupFilterContext, type GroupFilter } from "./GroupFilterContext";

export function GroupFilterProvider({ children }: { children: ReactNode }) {
  const [filter, setFilter] = useState<GroupFilter>({
    bow: null,
    identity: null,
    state: null,
  });

  return (
    <GroupFilterContext.Provider
      value={{
        filter,
        setFilter,
      }}
    >
      {children}
    </GroupFilterContext.Provider>
  );
}
