import { useState } from "react";
import { IconFileUpload } from "@tabler/icons-react";
import { useNavigate, useParams } from "react-router";
import { ActionIcon, rem, Stack } from "@mantine/core";
import {
  AddCompetitorModal,
  CompetitionTable,
  TableControlsContext,
  useAddCompetitorModalsStack,
  type TableControls,
} from "../../features";
import { TopBar } from "../../widgets";

export default function CompetitorsContent() {
  const { competitionId: paramCompetitionId } = useParams();
  const competitionId = Number(paramCompetitionId);

  const [tableControls, setTableControls] = useState<TableControls>({
    refresh: undefined,
  });

  const modalsStack = useAddCompetitorModalsStack();

  const navigate = useNavigate();

  const handleExcelTableUpload = () => {
    console.warn("handleExcelTableUpload temporary unavailable");
  };

  return (
    <>
      <AddCompetitorModal competitionId={competitionId} stack={modalsStack} />
      <Stack flex={1} style={{ overflow: "hidden" }} gap="lg" miw={rem(500)}>
        <TopBar
          title="Участники соревнования"
          onRefresh={tableControls.refresh}
          onAdd={() => modalsStack.open("add-competitor")}
          onBack={() => navigate("..")}
        >
          <ActionIcon onClick={handleExcelTableUpload}>
            <IconFileUpload />
          </ActionIcon>
        </TopBar>
        <TableControlsContext.Provider
          value={{
            controls: tableControls,
            setControls: setTableControls,
          }}
        >
          <CompetitionTable competitionId={competitionId} />
        </TableControlsContext.Provider>
      </Stack>
    </>
  );
}
