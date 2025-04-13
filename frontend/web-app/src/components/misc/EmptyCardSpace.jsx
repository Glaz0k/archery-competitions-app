import { Card, Center, Title } from "@mantine/core";

export default function EmptyCardSpace({ label }) {
  return (
    <Center flex={1}>
      <Card>
        <Title align="center" order={1}>
          {label}
        </Title>
      </Card>
    </Center>
  );
}
