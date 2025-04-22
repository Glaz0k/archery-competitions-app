import { Card, Center, Divider } from "@mantine/core";

export function ShotDivider() {
  return <Divider color="secondary.9" orientation="vertical" size="md" my="xs" />;
}

export function ShotCard({ children }) {
  return (
    <Card p="sm">
      <Center w={40} h={40}>
        {children}
      </Center>
    </Card>
  );
}
