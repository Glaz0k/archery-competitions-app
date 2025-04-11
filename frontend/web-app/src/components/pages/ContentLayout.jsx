import { Link, Outlet } from "react-router";
import {
  AppShell,
  Title,
  Text,
  Group,
  Button,
  ActionIcon,
  Flex,
  Box,
  Stack,
} from "@mantine/core";
import { IconLogout } from "@tabler/icons-react";

export default function ContentLayout() {
  return (
    <AppShell header={{ height: 100 }} padding={0}>
      <AppShell.Header>
        <Group h="100%">
          <Title order={1} flex={1}>
            ArcheryManager
          </Title>
          <Group>
            <Button component={Link} to={"/cups"}>
              <Title order={2}>Кубки</Title>
            </Button>
            <ActionIcon>
              <IconLogout />
            </ActionIcon>
          </Group>
        </Group>
      </AppShell.Header>
      <AppShell.Main style={{ display: "flex", flexDirection: "column" }}>
        <Box px={200} py={50} flex={1} display="flex" bg="gray">
          <Outlet />
        </Box>
        <Flex h={100} align="end" justify="center">
          <Text>©DevBow Team, {new Date().getFullYear()}</Text>
        </Flex>
      </AppShell.Main>
    </AppShell>
  );
}
