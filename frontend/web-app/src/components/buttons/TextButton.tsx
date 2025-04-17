import React, { forwardRef } from "react";
import { Button, ButtonProps, createPolymorphicComponent, Text } from "@mantine/core";

export interface TextButtonProps extends ButtonProps {
  label: string;
}

export const TextButton = createPolymorphicComponent<"button", TextButtonProps>(
  forwardRef<HTMLButtonElement, TextButtonProps>(({ label, size, children, ...others }, ref) => (
    <Button size={size} {...others} ref={ref}>
      <Text size={size}>{label}</Text>
    </Button>
  ))
);
