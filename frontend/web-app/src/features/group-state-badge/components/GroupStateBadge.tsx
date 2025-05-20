import { IconCheck, IconDots } from "@tabler/icons-react";
import { Badge, Text } from "@mantine/core";
import { GroupState } from "../../../entities";
import { getGroupStateDescription } from "../../../utils";

export interface GroupStateBadgeProps {
  state: GroupState;
}

export function GroupStateBadge({ state }: GroupStateBadgeProps) {
  const inner = <Text tt="capitalize">{getGroupStateDescription(state)}</Text>;
  switch (state) {
    case GroupState.COMPLETED:
      return (
        <Badge leftSection={<IconCheck />} color="green.8">
          {inner}
        </Badge>
      );
    default:
      return (
        <Badge leftSection={<IconDots />} color="dark.8">
          {inner}
        </Badge>
      );
  }
}
