import { useNavigate } from "react-router";
import { Card, Center, Stack, Title } from "@mantine/core";
import { TextButton } from "../widgets";

export default function NotFoundPage() {
  const navigate = useNavigate();
  return (
    <Center h="100vh">
      <Card>
        <Stack align="center">
          <Title order={1}>404</Title>
          <Title order={2}>Страница не найдена</Title>
          <TextButton label="Перейти на главную" onClick={() => navigate("/cups")} />
        </Stack>
      </Card>
    </Center>
  );
}
