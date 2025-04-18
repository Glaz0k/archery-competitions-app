import { IconCheck, IconDots } from "@tabler/icons-react";
import { Badge, Text } from "@mantine/core";
import GroupState from "../../enums/GroupState";

export default function GroupStateBadge({ state }) {
  switch (state?.value) {
    case GroupState.COMPLETED.value:
      return (
        <Badge leftSection={<IconCheck />} color="green.8">
          <Text tt="capitalize">{state?.textValue}</Text>
        </Badge>
      );
    default:
      return (
        <Badge leftSection={<IconDots />} color="dark.8">
          <Text tt="capitalize">{state?.textValue}</Text>
        </Badge>
      );
  }
}
