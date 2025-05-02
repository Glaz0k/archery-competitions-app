import { Select, type ComboboxData } from "@mantine/core";
import { Identity } from "../../../entities";
import { getIdentityDescription } from "../../../utils";
import { DEFAULT_ITEM } from "../../constants";

export function IdentitySelect({ setIdentity }: { setIdentity: (value: null | Identity) => void }) {
  const identityData: ComboboxData = Object.values(Identity).map((identityValue) => ({
    value: identityValue,
    label: getIdentityDescription(identityValue as Identity),
  }));

  return (
    <Select
      w="100%"
      label="Пол"
      data={identityData}
      clearable
      placeholder={DEFAULT_ITEM.label}
      onChange={(val) => setIdentity(val as Identity)}
    />
  );
}
