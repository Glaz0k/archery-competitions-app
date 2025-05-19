import { useContext } from "react";
import { TableControlsContext } from "./TableControlsContext";

export const useTableControls = () => {
  const context = useContext(TableControlsContext);

  if (!context) {
    throw new Error("useTableControls must be used within a TableControlsProvider");
  }

  return context;
};
