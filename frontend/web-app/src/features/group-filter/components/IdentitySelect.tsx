import { useEffect, useState } from "react";
import { NativeSelect, type ComboboxData } from "@mantine/core";
import { Identity } from "../../../entities";
import { getIdentityDescription } from "../../../utils";
import { DEFAULT_ITEM, DEFAULT_VALUE } from "../../constants";

export function IdentitySelect({
  setIdentity,
}: {
  setIdentity: (value: undefined | Identity) => void;
}) {
  const [selectedValue, setSelectedValue] = useState<string>(DEFAULT_VALUE);
  const identityValues = [DEFAULT_VALUE, ...Object.values(Identity)];

  const identityData: ComboboxData = identityValues.map((identityValue) => {
    if (identityValue === DEFAULT_VALUE) {
      return DEFAULT_ITEM;
    }
    return {
      value: identityValue,
      label: getIdentityDescription(identityValue as Identity),
    };
  });

  useEffect(() => {
    if (selectedValue === DEFAULT_VALUE) {
      setIdentity(undefined);
    } else {
      setIdentity(selectedValue as Identity);
    }
  }, [selectedValue, setIdentity]);

  return (
    <NativeSelect
      label="Пол"
      data={identityData}
      onChange={(e) => setSelectedValue(e.currentTarget.value)}
    />
  );
}
