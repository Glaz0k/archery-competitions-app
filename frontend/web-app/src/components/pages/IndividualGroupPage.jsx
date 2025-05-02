import { useEffect, useState } from "react";
import { useQueries } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router";
import { Stack, Tabs, Title } from "@mantine/core";
import { useDocumentTitle } from "@mantine/hooks";
import { getIndividualGroup } from "../../api/individualGroups";
import {
  COMPETITION_QUERY_KEYS,
  CUP_QUERY_KEYS,
  INDIVIDUAL_GROUP_QUERY_KEYS,
} from "../../api/queryKeys";
import NavigationBar from "../bars/NavigationBar";
import CompetitorsPanel from "./individual-group/CompetitorsPanel";
import FinalPanel from "./individual-group/FinalPanel";
import { GroupContext } from "./individual-group/GroupContext";
import QualificationPanel from "./individual-group/QualificationPanel";

export default function IndividualGroupPage() {
  const navigate = useNavigate();
  const { cupId, competitionId, groupId, groupSection } = useParams();

  const [info, setInfo] = useState({
    cup: null,
    competition: null,
    group: null,
  });
  const [controls, setControls] = useState({
    onExport: null,
    onRefresh: null,
    onEnd: null,
  });

  const [webTitle, setWebTitle] = useState(null);

  const mainQuery = useQueries({
    queries: [
      {
        queryKey: CUP_QUERY_KEYS.element(cupId),
        queryFn: () => getCup(cupId),
        initialData: null,
      },
      {
        queryKey: COMPETITION_QUERY_KEYS.element(competitionId),
        queryFn: () => getCompetition(competitionId),
        initialData: null,
      },
      {
        queryKey: INDIVIDUAL_GROUP_QUERY_KEYS.element(groupId),
        queryFn: () => getIndividualGroup(groupId),
        initialData: null,
      },
    ],
  });
  const [{ data: cup }, { data: competition }, { data: group }] = mainQuery;
  const [groupTitle, setGroupTitle] = useState("Дивизион");
  const [groupSubtitle, setGroupSubtitle] = useState("Кубок - Этап");

  const isMainInfoLoaded = cup && competition && group;
  const isMainInfoLoading = mainQuery.some((query) => query.isFetching);
  const isMainInfoLoadError = mainQuery.some((query) => query.isError);
  const isMainInfoNotFound = mainQuery.some((query) => query.error?.response?.status === 404);

  useDocumentTitle(webTitle);

  useEffect(() => {
    if (isMainInfoLoadError && isMainInfoNotFound) {
      navigate("/not-found");
    }
  }, [isMainInfoLoadError, isMainInfoNotFound, navigate]);

  useEffect(() => {
    if (isMainInfoLoaded) {
      setInfo({
        cup: cup,
        competition: competition,
        group: group,
      });
    }
  }, [isMainInfoLoaded, cup, competition, group]);

  useEffect(() => {
    if (competition && cup) {
      setGroupSubtitle(
        competition.stage.textValue +
          " - " +
          cup.title +
          (cup.season != null ? ", " + cup.season : "")
      );
    }
  }, [competition, cup]);

  useEffect(() => {
    if (group) {
      setGroupTitle(group.bow.textValue + " - " + group.identity.textValue);
    }
  }, [group]);

  useEffect(() => {
    let title = "ArcheryManager";
    if (cup && competition) {
      title += " - " + cup.title + " | " + competition.stage.textValue;
      if (competition.isEnded) {
        title += " - Завершён";
      }
      title += " | " + groupTitle;
    }
    setWebTitle(title);
  }, [cup, competition, groupTitle]);

  return (
    <Tabs
      keepMounted={false}
      variant="pills"
      value={groupSection}
      onChange={(value) => navigate(value)}
      display="flex"
      flex={1}
      style={{ overflow: "hidden" }}
    >
      <Stack gap="lg" display="flex" flex={1} style={{ overflow: "hidden" }}>
        <NavigationBar
          title={groupTitle}
          subTitle={groupSubtitle}
          loading={isMainInfoLoading}
          onBack={() => navigate("..")}
          onRefresh={controls?.onRefresh}
          onExport={controls?.onExport}
          onEnd={controls?.onEnd}
        >
          <Tabs.List>
            <Tabs.Tab value="competitors">
              <Title order={3}>{"Участники"}</Title>
            </Tabs.Tab>
            <Tabs.Tab value="qualification">
              <Title order={3}>{"Квалификация"}</Title>
            </Tabs.Tab>
            <Tabs.Tab value="final">
              <Title order={3}>{"Финал"}</Title>
            </Tabs.Tab>
          </Tabs.List>
        </NavigationBar>
        {isMainInfoLoaded && (
          <GroupContext.Provider
            value={{
              groupBow: group.bow,
            }}
          >
            <Tabs.Panel value="competitors">
              <CompetitorsPanel groupInfo={info} setGroupControls={setControls} />
            </Tabs.Panel>
            <Tabs.Panel value="qualification">
              <QualificationPanel groupInfo={info} setGroupControls={setControls} />
            </Tabs.Panel>
            <Tabs.Panel value="final" display="flex" flex={1} style={{ overflow: "hidden" }}>
              <FinalPanel groupInfo={info} setGroupControls={setControls} />
            </Tabs.Panel>
          </GroupContext.Provider>
        )}
      </Stack>
    </Tabs>
  );
}
