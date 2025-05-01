import { Stack } from "@mantine/core";
import { ControlsCard } from "../../../widgets";
import { useGroupFilter } from "../context/useGroupFilter";
import { BowSelect } from "./BowSelect";
import { IdentitySelect } from "./IdentitySelect";
import { StateSelect } from "./StateSelect";

export function GroupFilterControls() {
  const { filter, setFilter } = useGroupFilter();

  return (
    <ControlsCard>
      <Stack align="start" pos="relative" justify="stretch">
        <BowSelect
          setBow={(value) =>
            setFilter({
              ...filter,
              bow: value,
            })
          }
        />
        <IdentitySelect
          setIdentity={(value) =>
            setFilter({
              ...filter,
              identity: value,
            })
          }
        />
        <StateSelect
          setState={(value) =>
            setFilter({
              ...filter,
              state: value,
            })
          }
        />
      </Stack>
    </ControlsCard>
  );
}
