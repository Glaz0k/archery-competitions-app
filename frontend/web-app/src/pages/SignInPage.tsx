import { useNavigate } from "react-router";
import { Center, PasswordInput, Stack, Text, TextInput, Title } from "@mantine/core";
import { useForm, zodResolver } from "@mantine/form";
import { CredentialsSchema, type Credentials } from "../api";
import { useAdminSignIn } from "../api/endpoints/auth/hooks";
import { APP_NAME } from "../constants";
import { ControlsCard, TextButton } from "../widgets";

export default function SignInPage() {
  const signInForm = useForm<Credentials>({
    mode: "uncontrolled",
    initialValues: {
      login: "",
      password: "",
    },
    validate: zodResolver(CredentialsSchema),
  });
  const navigate = useNavigate();
  const {
    mutate: signIn,
    isPending: isSigningIn,
    isError: isSignInError,
  } = useAdminSignIn(() => navigate("/cups"));

  return (
    <Center h="100vh">
      <ControlsCard w={300}>
        <form onSubmit={signInForm.onSubmit((values) => signIn(values))}>
          <Stack align="center">
            <Title order={2}>{APP_NAME}</Title>
            <TextInput
              w="100%"
              label="Имя пользователя"
              placeholder="Username"
              key={signInForm.key("login")}
              {...signInForm.getInputProps("login")}
            />
            <PasswordInput
              w="100%"
              label="Пароль"
              placeholder="Password"
              key={signInForm.key("password")}
              {...signInForm.getInputProps("password")}
            />
            <TextButton type="submit" label="Войти" loading={isSigningIn} />
            {isSignInError && <Text c="red">Произошла ошибка</Text>}
          </Stack>
        </form>
      </ControlsCard>
    </Center>
  );
}
