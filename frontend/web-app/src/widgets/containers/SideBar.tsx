import type { PropsWithChildren } from "react";
import { rem, Stack } from "@mantine/core";

export function SideBar({ children }: PropsWithChildren) {
  return (
    <Stack w={rem(300)} p={0} gap="md" align="stretch">
      {children}
    </Stack>
  );
}
