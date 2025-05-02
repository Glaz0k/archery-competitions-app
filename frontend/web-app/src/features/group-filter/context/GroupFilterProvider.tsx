import { useState, type ReactNode } from "react";
import { GroupFilterContext, type GroupFilter } from "./GroupFilterContext";

export function GroupFilterProvider({ children }: { children: ReactNode }) {
  const [filter, setFilter] = useState<GroupFilter>({
    bow: undefined,
    identity: undefined,
    state: undefined,
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
