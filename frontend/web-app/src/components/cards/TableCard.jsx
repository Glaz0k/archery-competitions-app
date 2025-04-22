import { Card, LoadingOverlay, Table, useMantineTheme } from "@mantine/core";

export default function TableCard({ loading, children }) {
  const theme = useMantineTheme();
  return (
    <Card p={0} pos="relative">
      <LoadingOverlay visible={loading} />
      <Table.ScrollContainer>
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
