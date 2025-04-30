import { StrictMode } from "react";
import { setDefaultOptions } from "date-fns";
import { ru } from "date-fns/locale";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router";
import { MantineProvider } from "@mantine/core";
import { DatesProvider } from "@mantine/dates";
import { Notifications } from "@mantine/notifications";

import "dayjs/locale/ru";
import "@mantine/core/styles.css";
import "@mantine/dates/styles.css";
import "@mantine/notifications/styles.css";

import { QueryClientProvider } from "@tanstack/react-query";
import { queryClient } from "./api";
import App from "./App";
import { APP_NAME } from "./constants";
import theme from "./theme";

document.title = APP_NAME;

setDefaultOptions({ locale: ru });

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <QueryClientProvider client={queryClient}>
        <DatesProvider
          settings={{
            locale: "ru",
            firstDayOfWeek: 1,
            weekendDays: [0, 6],
          }}
        >
          <MantineProvider theme={theme}>
            <Notifications />
            <App />
          </MantineProvider>
        </DatesProvider>
      </QueryClientProvider>
    </BrowserRouter>
  </StrictMode>
);
