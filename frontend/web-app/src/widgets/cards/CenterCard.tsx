import { Card, Center, Title } from "@mantine/core";

export interface CenterCardProps {
  label: string;
}

export function CenterCard({ label }: CenterCardProps) {
  return (
    <Center flex={1}>
      <Card>
        <Title order={1}>{label}</Title>
      </Card>
    </Center>
  );
}
