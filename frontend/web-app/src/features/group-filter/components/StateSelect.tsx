import { useEffect, useState } from "react";
import { NativeSelect, type ComboboxData } from "@mantine/core";
import { GroupState } from "../../../entities";
import { getGroupStateDescription } from "../../../utils";
import { DEFAULT_ITEM, DEFAULT_VALUE } from "../../constants";

export function StateSelect({ setState }: { setState: (value: undefined | GroupState) => void }) {
  const [selectedValue, setSelectedValue] = useState<string>(DEFAULT_VALUE);
  const stateValues = [DEFAULT_VALUE, ...Object.values(GroupState)];

  const stateData: ComboboxData = stateValues.map((stateValue) => {
    if (stateValue === DEFAULT_VALUE) {
      return DEFAULT_ITEM;
    }
    return {
      value: stateValue,
      label: getGroupStateDescription(stateValue as GroupState),
    };
  });

  useEffect(() => {
    if (selectedValue === DEFAULT_VALUE) {
      setState(undefined);
    } else {
      setState(selectedValue as GroupState);
    }
  });

  return (
    <NativeSelect
      label="Состояние"
      data={stateData}
      onChange={(e) => setSelectedValue(e.currentTarget.value)}
    />
  );
}
