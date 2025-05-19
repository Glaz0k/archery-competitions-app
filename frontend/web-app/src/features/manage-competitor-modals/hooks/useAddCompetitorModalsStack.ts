import { useModalsStack } from "@mantine/core";

export type AddCompetitorModalsStackType = "add-competitor" | "register-competitor";

export const addCompetitorModalsStack: AddCompetitorModalsStackType[] = [
  "add-competitor",
  "register-competitor",
];

export const useAddCompetitorModalsStack = () => {
  return useModalsStack<AddCompetitorModalsStackType>(addCompetitorModalsStack);
};
