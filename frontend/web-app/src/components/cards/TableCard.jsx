import { Card, LoadingOverlay, Table } from "@mantine/core";

export default function TableCard({ loading, children }) {
  return (
    <Card p={0} pos="relative">
      <LoadingOverlay visible={loading} />
      <Table.ScrollContainer>
        <Table tabularNums withColumnBorders>
          {children}
        </Table>
      </Table.ScrollContainer>
    </Card>
  );
}
