import { Select, type ComboboxData } from "@mantine/core";
import { GroupState } from "../../../entities";
import { getGroupStateDescription } from "../../../utils";
import { DEFAULT_ITEM } from "../../constants";

export function StateSelect({ setState }: { setState: (value: null | GroupState) => void }) {
  const stateData: ComboboxData = Object.values(GroupState).map((stateValue) => ({
    value: stateValue,
    label: getGroupStateDescription(stateValue as GroupState),
  }));

  return (
    <Select
      w="100%"
      label="Состояние"
      data={stateData}
      clearable
      placeholder={DEFAULT_ITEM.label}
      onChange={(val) => setState(val as GroupState)}
    />
  );
}
