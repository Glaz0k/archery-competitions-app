import { createContext } from "react";

export const FinalContext = createContext({
  selectedPlaceId: null,
  setSelectedPlaceId: undefined,
  sparringSize: {
    heigth: null,
    width: null,
  },
});
