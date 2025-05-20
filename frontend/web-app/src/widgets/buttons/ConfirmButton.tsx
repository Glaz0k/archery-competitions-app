import { forwardRef } from "react";
import { IconCheck } from "@tabler/icons-react";
import { createPolymorphicComponent } from "@mantine/core";
import { TextButton, type TextButtonProps } from "./TextButton";

export const ConfirmButton = createPolymorphicComponent<"button", TextButtonProps>(
  forwardRef<HTMLButtonElement, TextButtonProps>(({ ...props }, ref) => (
    <TextButton
      component="button"
      ref={ref}
      variant="filled"
      color="green.8"
      rightSection={<IconCheck />}
      {...props}
    />
  ))
);
