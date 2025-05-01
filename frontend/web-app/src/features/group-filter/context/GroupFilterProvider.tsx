import { useMemo, useState, type ReactNode } from "react";
import { GroupFilterContext, type GroupFilter } from "./GroupFilterContext";

const initialState: GroupFilter = {
  bow: undefined,
  identity: undefined,
  state: undefined,
};

export function GroupFilterProvider({ children }: { children: ReactNode }) {
  const [filter, setFilter] = useState(initialState);

  const value = useMemo(
    () => ({
      filter,
      setFilter,
    }),
    [filter]
  );

  return <GroupFilterContext.Provider value={value}>{children}</GroupFilterContext.Provider>;
}
