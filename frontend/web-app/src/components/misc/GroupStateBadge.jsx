import { IconCheck, IconDots } from "@tabler/icons-react";
import { Badge } from "@mantine/core";
import GroupState from "../../enums/GroupState";

export default function GroupStateBadge({ state }) {
  switch (state) {
    case GroupState.COMPLETED:
      return <Badge leftSection={<IconCheck />}>{state.textValue}</Badge>;
    default:
      return <Badge leftSection={<IconDots />}>{state.textValue}</Badge>;
  }
}
