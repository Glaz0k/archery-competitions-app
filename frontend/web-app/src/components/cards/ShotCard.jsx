import { AspectRatio, Card, Center, Title } from "@mantine/core";

export default function ShotCard({ score, editing, children }) {
  return (
    <AspectRatio ratio={1}>
      <Card p="sm">
        <Center w={40} h={40}>
          {editing ? children : <Title order={2}>{score || "-"}</Title>}
        </Center>
      </Card>
    </AspectRatio>
  );
}
