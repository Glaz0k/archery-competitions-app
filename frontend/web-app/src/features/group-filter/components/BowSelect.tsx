import { Select, type ComboboxData } from "@mantine/core";
import { BowClass } from "../../../entities";
import { getBowClassDescription } from "../../../utils";
import { DEFAULT_ITEM } from "../../constants";

export function BowSelect({ setBow }: { setBow: (value: null | BowClass) => void }) {
  const bowData: ComboboxData = Object.values(BowClass).map((bowValue) => ({
    value: bowValue,
    label: getBowClassDescription(bowValue),
  }));

  return (
    <Select
      w="100%"
      label="Класс лука"
      data={bowData}
      clearable
      placeholder={DEFAULT_ITEM.label}
      onChange={(val) => setBow(val as BowClass)}
    />
  );
}
