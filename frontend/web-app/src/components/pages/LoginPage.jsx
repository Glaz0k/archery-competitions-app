import {
  Button,
  Center,
  Paper,
  PasswordInput,
  Stack,
  TextInput,
  Title,
} from "@mantine/core";
import { useForm } from "@mantine/form";

export default function LoginPage() {
  const form = useForm({
    mode: "uncontrolled",
    initialValues: {
      username: "",
      password: "",
    },
  });

  return (
    <Center style={{ height: "100vh" }}>
      <Paper>
        <form onSubmit={form.onSubmit((values) => console.log(values))}>
          <Stack align="stretch" w="300px">
            <Title align="center" order={2}>
              ArcheryManager
            </Title>
            <TextInput
              label="Имя пользователя"
              placeholder="Username"
              key={form.key("username")}
              {...form.getInputProps("username")}
            />
            <PasswordInput
              label="Пароль"
              placeholder="Password"
              key={form.key("password")}
              {...form.getInputProps("password")}
            />
            <Button type="submit">Войти</Button>
          </Stack>
        </form>
      </Paper>
    </Center>
  );
}
