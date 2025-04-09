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

  const handleSubmit = function (values) {
    console.log(values);
  };

  return (
    <Center style={{ height: "100vh" }}>
      <Paper shadow="md" radius="md">
        <form onSubmit={form.onSubmit(handleSubmit)}>
          <Stack align="center" w={300}>
            <Title order={2}>ArcheryManager</Title>
            <TextInput
              label="Имя пользователя"
              placeholder="Username"
              key={form.key("username")}
              {...form.getInputProps("username")}
              w="100%"
            />
            <PasswordInput
              label="Пароль"
              placeholder="Password"
              key={form.key("password")}
              {...form.getInputProps("password")}
              w="100%"
            />
            <Button type="submit">Войти</Button>
          </Stack>
        </form>
      </Paper>
    </Center>
  );
}
