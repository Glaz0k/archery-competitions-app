import { createTheme, MantineProvider } from "@mantine/core";
import LoginPage from "./components/pages/LoginPage";

const theme = createTheme();

export default function App() {
  return (
    <MantineProvider theme={theme}>
      <LoginPage />
    </MantineProvider>
  );
}
