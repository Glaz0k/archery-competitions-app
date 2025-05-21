import { forwardRef } from "react";
import { IconX } from "@tabler/icons-react";
import { createPolymorphicComponent } from "@mantine/core";
import { TextButton, type TextButtonProps } from "./TextButton";

export const CancelButton = createPolymorphicComponent<"button", TextButtonProps>(
  forwardRef<HTMLButtonElement, TextButtonProps>(({ ...props }, ref) => (
    <TextButton
      component="button"
      ref={ref}
      variant="filled"
      color="dark.3"
      rightSection={<IconX />}
      {...props}
    />
  ))
);
