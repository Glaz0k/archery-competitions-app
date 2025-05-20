import { Stack, Title } from "@mantine/core";
import { ControlsCard } from "../../../widgets";
import { useGroupFilter } from "../context/useGroupFilter";
import { BowSelect } from "./BowSelect";
import { IdentitySelect } from "./IdentitySelect";
import { StateSelect } from "./StateSelect";

export function GroupFilterControls() {
  const { setFilter } = useGroupFilter();

  return (
    <ControlsCard>
      <Title order={3}>Фильтры</Title>
      <Stack align="start" pos="relative" justify="stretch" gap="sm">
        <BowSelect
          setBow={(bow) =>
            setFilter((prev) => ({
              ...prev,
              bow,
            }))
          }
        />
        <IdentitySelect
          setIdentity={(identity) =>
            setFilter((prev) => ({
              ...prev,
              identity,
            }))
          }
        />
        <StateSelect
          setState={(state) =>
            setFilter((prev) => ({
              ...prev,
              state,
            }))
          }
        />
      </Stack>
    </ControlsCard>
  );
}
