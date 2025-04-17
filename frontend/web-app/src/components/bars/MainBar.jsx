import { IconArrowLeft, IconPlus, IconRefresh } from "@tabler/icons-react";
import { ActionIcon, Group, Title } from "@mantine/core";
import PrimaryCard from "../cards/PrimaryCard";

export default function MainBar({ title, onRefresh, onAdd, onBack, children }) {
  return (
    <PrimaryCard>
      <Group>
        {onBack && (
          <ActionIcon onClick={onBack}>
            <IconArrowLeft />
          </ActionIcon>
        )}
        <Title order={2} flex={1}>
          {title}
        </Title>
        {children}
        <ActionIcon onClick={onRefresh}>
          <IconRefresh />
        </ActionIcon>
        <ActionIcon onClick={onAdd}>
          <IconPlus />
        </ActionIcon>
      </Group>
    </PrimaryCard>
  );
}
