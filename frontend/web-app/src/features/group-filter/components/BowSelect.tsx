import { useEffect, useState } from "react";
import { NativeSelect, type ComboboxData } from "@mantine/core";
import { BowClass } from "../../../entities";
import { getBowClassDescription } from "../../../utils";
import { DEFAULT_ITEM, DEFAULT_VALUE } from "../../constants";

export function BowSelect({ setBow }: { setBow: (value: undefined | BowClass) => void }) {
  const [selectedValue, setSelectedValue] = useState<string>(DEFAULT_VALUE);
  const bowValues = [DEFAULT_VALUE, ...Object.values(BowClass)];

  const bowData: ComboboxData = bowValues.map((bowValue) => {
    if (bowValue === DEFAULT_VALUE) {
      return DEFAULT_ITEM;
    }
    return {
      value: bowValue,
      label: getBowClassDescription(bowValue as BowClass),
    };
  });

  useEffect(() => {
    if (selectedValue === DEFAULT_VALUE) {
      setBow(undefined);
    } else {
      setBow(selectedValue as BowClass);
    }
  }, [selectedValue, setBow]);

  return (
    <NativeSelect
      label="Класс лука"
      data={bowData}
      onChange={(e) => setSelectedValue(e.currentTarget.value)}
    />
  );
}
