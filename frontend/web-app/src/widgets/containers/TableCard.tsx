import type { PropsWithChildren } from "react";
import { Card, LoadingOverlay, Table, useMantineTheme } from "@mantine/core";

export interface TableCardProps {
  loading: boolean;
}

export function TableCard({ loading = false, children }: PropsWithChildren<TableCardProps>) {
  const theme = useMantineTheme();
  return (
    <Card p={0} pos="relative">
      <LoadingOverlay visible={loading} />
      <Table.ScrollContainer minWidth={500}>
        <Table
          tabularNums
          withColumnBorders
          highlightOnHover
          highlightOnHoverColor={`${theme.colors.gray[0]}33`}
        >
          {children}
        </Table>
      </Table.ScrollContainer>
    </Card>
  );
}
