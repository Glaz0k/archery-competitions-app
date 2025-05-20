import { useContext } from "react";
import { GroupSectionsContext } from "./GroupSectionsContext";

export const useGroupSections = () => {
  const context = useContext(GroupSectionsContext);

  if (!context) {
    throw new Error("useTableControls must be used within a TableControlsProvider");
  }

  return context;
};
