import {
  ActionIcon,
  Anchor,
  AppShellHeader,
  Badge,
  Button,
  Card,
  colorsTuple,
  createTheme,
  Input,
  InputLabel,
  InputWrapper,
  Modal,
  ModalBase,
  ModalBaseBody,
  Pagination,
  Paper,
  rem,
  Skeleton,
  Text,
  TextInput,
  ThemeIcon,
} from "@mantine/core";
import { DatePickerInput } from "@mantine/dates";

const light = [
  "#f5f5f5",
  "#e7e7e7",
  "#cdcdcd",
  "#b2b2b2",
  "#9a9a9a",
  "#8b8b8b",
  "#848484",
  "#717171",
  "#656565",
  "#575757",
];

const primary = [
  "#f0faf1",
  "#e0f1e1",
  "#bce3be",
  "#95d498",
  "#75c778",
  "#61bf63",
  "#56bc58",
  "#46a548",
  "#3c933f",
  "#2e7d32",
];

const secondary = [
  "#f3f5f7",
  "#e8e8e8",
  "#cdcfd0",
  "#aeb4b8",
  "#949ea4",
  "#839098",
  "#798a93",
  "#667780",
  "#596a73",
  "#37474f",
];

const white = colorsTuple("#FFFFFF");

const theme = createTheme({
  primaryColor: "primary",
  autoContrast: true,
  colors: {
    light,
    primary,
    secondary,
    white,
  },
  fontFamily: "Inter, sans-serif",
  fontSizes: {
    xs: rem(12),
    sm: rem(14),
    md: rem(16),
    lg: rem(18),
    xl: rem(20),
  },
  headings: {
    fontFamily: "Manrope, sans-serif",
    sizes: {
      h1: {
        fontSize: rem(36),
        fontWeight: 600,
      },
      h2: {
        fontSize: rem(28),
        fontWeight: 500,
      },
      h3: {
        fontSize: rem(22),
        fontWeight: 400,
      },
    },
  },
  defaultRadius: 8,
  spacing: {
    xs: rem(8),
    sm: rem(12),
    md: rem(20),
    lg: rem(50),
    xl: rem(200),
  },
  components: {
    AppShellHeader: AppShellHeader.extend({
      defaultProps: {
        bg: "primary.9",
        c: "white.0",
        p: "md",
        bd: "none",
      },
    }),
    Text: Text.extend({
      styles: (_theme, params) => {
        let style = { root: {} };
        switch (params.fz) {
          case "xs":
            style.root = { fontWeight: 200 };
            break;
          case "sm":
            style.root = { fontWeight: 300 };
            break;
          case "md":
            style.root = { fontWeight: 400 };
            break;
        }
        return style;
      },
    }),
    Button: Button.extend({
      defaultProps: {
        size: "sm",
        variant: "light",
        color: "white.0",
      },
    }),
    ActionIcon: ActionIcon.extend({
      defaultProps: {
        size: "xl",
        radius: "xl",
        variant: "outline",
        color: "white.0",
      },
    }),
    ThemeIcon: ThemeIcon.extend({
      defaultProps: {
        size: "xl",
        radius: "xl",
        variant: "outline",
        bd: "none",
        color: "white.0",
      },
    }),
    Card: Card.extend({
      defaultProps: {
        bg: "secondary.9",
        c: "white.0",
        p: "md",
      },
    }),
    Pagination: Pagination.extend({
      defaultProps: {
        color: "secondary.9",
      },
    }),
    Anchor: Anchor.extend({
      styles: (_theme) => ({
        root: {
          color: "inherit",
        },
      }),
    }),
    Modal: Modal.extend({
      defaultProps: {
        centered: true,
        lockScroll: false,
        padding: "md",
      },
    }),
    DatePickerInput: DatePickerInput.extend({
      defaultProps: {
        valueFormat: "D MMMM YYYY",
      },
    }),
    Skeleton: Skeleton.extend({
      defaultProps: {
        opacity: 0.7,
      },
    }),
    Badge: Badge.extend({
      defaultProps: {
        size: "xl",
        variant: "white",
        color: "green.8",
      },
    }),
  },
});

export default theme;
