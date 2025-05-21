import type { PropsWithChildren } from "react";
import { IconCircleFilled } from "@tabler/icons-react";
import { Card, Group, LoadingOverlay, ThemeIcon, Title, useMantineTheme } from "@mantine/core";

export interface RangeCardProps {
  active: boolean;
  title: string;
  loading: boolean;
}

export function RangeCard({ active, title, loading, children }: PropsWithChildren<RangeCardProps>) {
  const theme = useMantineTheme();
  return (
    <Card bg="gray.0" bd="5px solid #E0E0E0" c={theme.black}>
      <LoadingOverlay visible={loading} />
      <Group gap="md">
        <ThemeIcon color={active ? theme.colors.green[6] : theme.colors.gray[6]}>
          <IconCircleFilled />
        </ThemeIcon>
        <Title order={2} flex={1}>
          {title}
        </Title>
        {children}
      </Group>
    </Card>
  );
}
