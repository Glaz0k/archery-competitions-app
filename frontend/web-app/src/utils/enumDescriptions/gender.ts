import { Gender } from "../../entities";

export const genderDescriptions: Record<Gender, string> = {
  [Gender.MALE]: "Мужчина",
  [Gender.FEMALE]: "Женщина",
};

export const getGenderDescription = (gender: Gender): string => {
  return genderDescriptions[gender] || "Неизвестно";
};
