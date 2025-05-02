import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Tabs, Title } from "@mantine/core";
import { getSection } from "../../api/qualificationSections";
import { SECTION_QUERY_KEYS } from "../../api/queryKeys";
import NavigationBar from "../bars/NavigationBar";
import QualificationRangeGroup from "./QualificationRangeGroup";

export default function QualificationSection({ sectionId }) {
  const [selectedRoundOrdinal, setSelectedRoundOrdinal] = useState(null);

  const [refreshFn, setRefreshFn] = useState(null);

  const { data: section, isFetching: isSectionFetching } = useQuery({
    queryKey: SECTION_QUERY_KEYS.element(sectionId),
    queryFn: () => getSection(sectionId),
    initialData: null,
  });

  const isSectionLoading = isSectionFetching;

  const roundTabs = [...(section?.rounds || [])]
    .sort((a, b) => a.roundOrdinal - b.roundOrdinal)
    .map((round, index) => (
      <Tabs.Tab key={index} value={String(round.roundOrdinal)}>
        <Title order={3}>{"Раунд " + round.roundOrdinal} </Title>
      </Tabs.Tab>
    ));

  if (section == null && !isSectionLoading) {
    return <NotFoundCard label={"Произошла ошибка"} />;
  }

  return (
    <>
      <NavigationBar
        title={section?.competitor.fullName || "Загрузка..."}
        loading={isSectionLoading}
        onRefresh={refreshFn}
      >
        <Tabs
          variant="pills"
          value={selectedRoundOrdinal}
          onChange={(value) => setSelectedRoundOrdinal(value)}
        >
          <Tabs.List>{roundTabs}</Tabs.List>
        </Tabs>
      </NavigationBar>
      {selectedRoundOrdinal != null && (
        <QualificationRangeGroup
          sectionId={section?.id}
          roundOrdinal={selectedRoundOrdinal}
          setRefresh={setRefreshFn}
        />
      )}
    </>
  );
}
