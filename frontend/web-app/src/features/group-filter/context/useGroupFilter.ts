import { useContext } from "react";
import { GroupFilterContext } from "./GroupFilterContext";

export const useGroupFilter = () => {
  const context = useContext(GroupFilterContext);

  if (!context) {
    throw new Error("useGroupFilter must be used within a GroupFilterProvider");
  }

  return context;
};
