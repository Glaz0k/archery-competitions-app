import { IconLogout } from "@tabler/icons-react";
import { Link, Outlet } from "react-router";
import { ActionIcon, AppShell, Box, Button, Flex, Group, Text, Title } from "@mantine/core";

export default function ContentLayout() {
  return (
    <AppShell header={{ height: 100 }} padding={0}>
      <AppShell.Header>
        <Group h="100%">
          <Title order={1} flex={1}>
            ArcheryManager
          </Title>
          <Group>
            <Button size="lg" component={Link} to={"/cups"}>
              <Title order={2}>{"Кубки"}</Title>
            </Button>
            <ActionIcon>
              <IconLogout />
            </ActionIcon>
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
