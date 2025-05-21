import { forwardRef } from "react";
import { Button, createPolymorphicComponent, Text, type ButtonProps } from "@mantine/core";

export interface TextButtonProps extends ButtonProps {
  label: string;
}

export const TextButton = createPolymorphicComponent<"button", TextButtonProps>(
  forwardRef<HTMLButtonElement, TextButtonProps>(({ label, size, ...others }, ref) => (
    <Button component="button" size={size} ref={ref} {...others}>
      <Text size={size}>{label}</Text>
    </Button>
  ))
);
