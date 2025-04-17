import React from "react";
import { Card, CardProps } from "@mantine/core";

export default function PrimaryCard({ children, ...props }: CardProps) {
  return (
    <Card {...props} bg="primary.9" color="white.0">
      {children}
    </Card>
  );
}
