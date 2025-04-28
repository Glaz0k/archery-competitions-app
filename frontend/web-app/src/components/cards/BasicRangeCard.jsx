import { IconCircleFilled } from "@tabler/icons-react";
import { Card, Group, LoadingOverlay, ThemeIcon, Title, useMantineTheme } from "@mantine/core";

export default function BasicRangeCard({ active, title, loading, children }) {
  const theme = useMantineTheme();

  return (
    <Card bg="gray.0" bd="5px solid #E0E0E0" c={theme.black}>
      <LoadingOverlay visible={loading} />
      <Group gap="md">
        <ThemeIcon
          color={
            active == null
              ? theme.colors.gray[6]
              : active
                ? theme.colors.green[6]
                : theme.colors.yellow[6]
          }
        >
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
