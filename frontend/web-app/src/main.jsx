import { StrictMode } from "react";
import { setDefaultOptions } from "date-fns";
import { ru } from "date-fns/locale";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router";
import { createTheme, MantineProvider } from "@mantine/core";
import { DatePickerInput, DatesProvider } from "@mantine/dates";

import "dayjs/locale/ru";

import App from "./App.jsx";

import "@mantine/core/styles.css";
import "@mantine/dates/styles.css";

const theme = createTheme({
  components: {
    DatePickerInput: {
      defaultProps: {
        valueFormat: "D MMMM YYYY",
      },
    },
  },
});

setDefaultOptions({ locale: ru });

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <BrowserRouter>
      <DatesProvider
        settings={{
          locale: "ru",
          firstDayOfWeek: 1,
          weekendDays: [0, 6],
        }}
      >
        <MantineProvider theme={theme}>
          <App />
        </MantineProvider>
      </DatesProvider>
    </BrowserRouter>
  </StrictMode>
);
