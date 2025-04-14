import { IconPlus } from "@tabler/icons-react";
import { ActionIcon, Group, Title } from "@mantine/core";

export default function MainBar({ title, onAdd, children }) {
  return (
    <Group>
      <Title order={2} flex={1}>
        {title}
      </Title>
      {children}
      <ActionIcon onClick={onAdd}>
        <IconPlus />
      </ActionIcon>
    </Group>
  );
}
