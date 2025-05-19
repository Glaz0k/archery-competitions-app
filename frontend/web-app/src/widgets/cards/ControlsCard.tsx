import { Card, type CardProps } from "@mantine/core";

export function ControlsCard({ children, ...props }: CardProps) {
  return (
    <Card bg="primary.9" color="white.0" {...props}>
      {children}
    </Card>
  );
}
