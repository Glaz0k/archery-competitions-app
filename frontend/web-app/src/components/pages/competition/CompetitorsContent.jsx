import { IconFileUpload } from "@tabler/icons-react";
import { useNavigate, useParams } from "react-router";
import { ActionIcon, rem, Stack } from "@mantine/core";
import { useCompetitionCompetitors } from "../../../api";
import { CompetitionTable } from "../../../features";
import { TopBar } from "../../../widgets";

export default function CompetitorsContent() {
  const { competitionId: paramCompetitionId } = useParams();
  const competitionId = Number(paramCompetitionId);

  const { refetch: refetchCompetitorDetails } = useCompetitionCompetitors(competitionId);

  const navigate = useNavigate();

  const handleCompetitorAdd = () => {
    // TODO
    console.warn("handleCompetitorAdd temporary unavailable");
  };

  const handleExcelTableUpload = () => {
    // TODO
    console.warn("handleExcelTableUpload temporary unavailable");
  };

  return (
    <>
      <Stack flex={1} style={{ overflow: "hidden" }} gap="lg" miw={rem(500)}>
        <TopBar
          title="Участники соревнования"
          onRefresh={refetchCompetitorDetails}
          onAdd={handleCompetitorAdd}
          onBack={() => navigate("..")}
        >
          <ActionIcon onClick={handleExcelTableUpload}>
            <IconFileUpload />
          </ActionIcon>
        </TopBar>
        <CompetitionTable competitionId={competitionId} />
      </Stack>
    </>
  );
}
