import { Link, Outlet } from "react-router";
import {
  AppShell,
  Title,
  Text,
  Group,
  Button,
  ActionIcon,
  Flex,
} from "@mantine/core";
import { IconLogout } from "@tabler/icons-react";

export default function ContentLayout() {
  return (
    <AppShell header={{ height: 100 }} footer={{ height: 100 }}>
      <AppShell.Header>
        <Flex align="center" h="100%">
          <Title order={1} flex={1}>
            ArcheryManager
          </Title>
          <Group>
            <Button component={Link} to={"/cups"}>
              Кубки
            </Button>
            <ActionIcon>
              <IconLogout />
            </ActionIcon>
          </Group>
        </Flex>
      </AppShell.Header>
      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
      <AppShell.Footer>
        <Flex align="end" h="100%" justify="center">
          <Text>©DevBow Team, {new Date().getFullYear()}</Text>
        </Flex>
      </AppShell.Footer>
    </AppShell>
  );
}
