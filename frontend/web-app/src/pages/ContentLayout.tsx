import { IconHome, IconLogout, IconUsers } from "@tabler/icons-react";
import { Link, Outlet, useNavigate } from "react-router";
import { ActionIcon, AppShell, Box, Flex, Group, Text, Title, Tooltip } from "@mantine/core";
import { useLogout } from "../api";
import { APP_NAME } from "../constants";

export default function ContentLayout() {
  const navigate = useNavigate();
  const { mutate: logout, isPending: isLoggingOut } = useLogout(() => navigate("/sign-in"));
  return (
    <AppShell header={{ height: 100 }} padding={0}>
      <AppShell.Header>
        <Group h="100%">
          <Title order={1} flex={1}>
            {APP_NAME}
          </Title>
          <Group>
            <Tooltip label="Домашняя - Кубки">
              <ActionIcon component={Link} to="/cups">
                <IconHome />
              </ActionIcon>
            </Tooltip>
            <Tooltip label="Пользователи">
              <ActionIcon component={Link} to="/competitors">
                <IconUsers />
              </ActionIcon>
            </Tooltip>
            <Tooltip label="Выйти">
              <ActionIcon onClick={() => logout()} loading={isLoggingOut}>
                <IconLogout />
              </ActionIcon>
            </Tooltip>
          </Group>
        </Group>
      </AppShell.Header>
      <AppShell.Main style={{ display: "flex", flexDirection: "column" }}>
        <Box px="xl" py="lg" flex={1} display="flex" pos="relative">
          <Outlet />
        </Box>
        <Flex h={100} align="end" justify="center" bg="secondary.9" c="white.0" p="md">
          <Text fz="sm">
            {"©DevBow Team, "}
            {new Date().getFullYear()}
          </Text>
        </Flex>
      </AppShell.Main>
    </AppShell>
  );
}
