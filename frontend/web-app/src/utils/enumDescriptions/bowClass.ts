import { BowClass } from "../../entities";

export const bowClassDescriptions: Record<BowClass, string> = {
  [BowClass.CLASSIC]: "Классический лук",
  [BowClass.BLOCK]: "Блочный лук",
  [BowClass.CLASSIC_NEWBIE]: "Классический лук (новички)",
  [BowClass.CLASSIC_3D]: "3Д-классический лук",
  [BowClass.COMPOUND_3D]: "3Д-составной лук",
  [BowClass.LONG_3D]: "3Д-длинный лук",
  [BowClass.PERIPHERAL]: "Периферийный лук",
  [BowClass.PERIPHERAL_WITH_RING]: "Периферийный лук (с кольцом)",
};

export const getBowClassDescription = (bow: BowClass): string => {
  return bowClassDescriptions[bow] || "Неизвестно";
};
