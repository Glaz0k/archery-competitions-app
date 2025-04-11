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
  const loginForm = useForm({
    mode: "uncontrolled",
    initialValues: {
      username: "",
      password: "",
    },
  });

  const handleLoginSubmit = (loginFormValues) => {
    console.log(loginFormValues);
  };

  return (
    <Center style={{ height: "100vh" }}>
      <Paper shadow="md" radius="md">
        <form onSubmit={loginForm.onSubmit(handleLoginSubmit)}>
          <Stack align="center" w={300}>
            <Title order={2}>ArcheryManager</Title>
            <TextInput
              label="Имя пользователя"
              placeholder="Username"
              key={loginForm.key("username")}
              {...loginForm.getInputProps("username")}
              w="100%"
            />
            <PasswordInput
              label="Пароль"
              placeholder="Password"
              key={loginForm.key("password")}
              {...loginForm.getInputProps("password")}
              w="100%"
            />
            <Button type="submit">Войти</Button>
          </Stack>
        </form>
      </Paper>
    </Center>
  );
}
