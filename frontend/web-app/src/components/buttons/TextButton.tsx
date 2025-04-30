import { Button, Text, type ButtonProps, type ElementProps } from "@mantine/core";

export interface TextButtonProps extends ButtonProps, ElementProps<"button", keyof ButtonProps> {
  label: string;
}

export function TextButton({ label, size, ...others }: TextButtonProps) {
  return (
    <Button size={size} {...others}>
      <Text size={size}>{label}</Text>
    </Button>
  );
}
