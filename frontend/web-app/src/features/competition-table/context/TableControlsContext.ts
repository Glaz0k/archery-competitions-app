import { createContext, type Dispatch, type SetStateAction } from "react";

export interface TableControls {
  refresh: (() => void) | undefined;
}

export interface TableControlsContextType {
  controls: TableControls;
  setControls: Dispatch<SetStateAction<TableControls>>;
}

export const TableControlsContext = createContext<TableControlsContextType | null>(null);
