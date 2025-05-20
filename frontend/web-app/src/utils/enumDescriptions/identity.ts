import { Identity } from "../../entities";

export const identityDescriptions: Record<Identity, string> = {
  [Identity.MALES]: "Мужчины",
  [Identity.FEMALES]: "Женщины",
  [Identity.UNITED]: "Объединенный",
};

export const getIdentityDescription = (identity: Identity): string => {
  return identityDescriptions[identity] || "Неизвестно";
};
