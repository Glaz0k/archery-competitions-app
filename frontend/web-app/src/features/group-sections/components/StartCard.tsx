import { IconFlag } from "@tabler/icons-react";
import { Center, Stack, Title } from "@mantine/core";
import { ControlsCard, TextButton } from "../../../widgets";

export interface StartCardProps {
  title: string;
  loading?: boolean;
  onStart: () => unknown;
}

export function StartCard({ title, loading = true, onStart }: StartCardProps) {
  return (
    <Center flex={1}>
      <ControlsCard>
        <Stack>
          <Title order={3}>{title}</Title>
          <TextButton
            label="Начать"
            rightSection={<IconFlag />}
            loading={loading}
            onClick={onStart}
          />
        </Stack>
      </ControlsCard>
    </Center>
  );
}
