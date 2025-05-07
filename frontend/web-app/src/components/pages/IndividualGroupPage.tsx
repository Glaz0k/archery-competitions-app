import { useEffect, useMemo, useState } from "react";
import { isAxiosError } from "axios";
import { useNavigate, useParams } from "react-router";
import { Stack, Tabs, Text, Title } from "@mantine/core";
import { useDocumentTitle } from "@mantine/hooks";
import { useCompetition, useCup, useIndividualGroup } from "../../api";
import { APP_NAME } from "../../constants";
import {
  CompetitorsSection,
  GroupSectionsContext,
  QualificationSection,
  type GroupSections,
} from "../../features";
import {
  getBowClassDescription,
  getCompetitionStageDescription,
  getIdentityDescription,
} from "../../utils";
import { PageLoader, TabsBar } from "../../widgets";

export default function IndividualGroupPage() {
  const navigate = useNavigate();
  const {
    cupId: paramCupId,
    competitionId: paramCompetitionId,
    groupId: paramGroupId,
    groupSection,
  } = useParams();
  const cupId = Number(paramCupId);
  const competitionId = Number(paramCompetitionId);
  const groupId = Number(paramGroupId);

  const [webTitle, setWebTitle] = useState<string>("");
  useDocumentTitle(webTitle);

  const [groupSections, setGroupSections] = useState<GroupSections>({
    info: {
      cup: null,
      competition: null,
      group: null,
    },
    controls: {
      onRefresh: undefined,
      onComplete: undefined,
      onExport: undefined,
    },
  });

  const {
    data: cup,
    isFetching: isCupFetching,
    isLoading: isCupLoading,
    isError: isCupError,
    error: cupError,
  } = useCup(cupId);
  const {
    data: competition,
    isFetching: isCompetitionFetching,
    isLoading: isCompetitionLoading,
    isError: isCompetitionError,
    error: competitionError,
  } = useCompetition(competitionId);
  const {
    data: group,
    isFetching: isGroupFetching,
    isLoading: isGroupLoading,
    isError: isGroupError,
    error: groupError,
  } = useIndividualGroup(groupId);

  const isMainInfoLoading = isCupLoading || isCompetitionLoading || isGroupLoading;
  const isMainInfoFetching = isCupFetching || isCompetitionFetching || isGroupFetching;
  const isMainInfoError = isCupError || isCompetitionError || isGroupError;
  const mainInfoErrors = useMemo(
    () => [cupError, competitionError, groupError],
    [competitionError, cupError, groupError]
  );

  useEffect(() => {
    if (cup && competition && group) {
      setGroupSections((prev) => ({
        ...prev,
        info: {
          cup: cup,
          competition: competition,
          group: group,
        },
      }));
    }
  }, [competition, cup, group]);

  useEffect(() => {
    if (isMainInfoError && mainInfoErrors.some((e) => isAxiosError(e) && e.status === 404)) {
      navigate("/not-found");
    }
  }, [isMainInfoError, mainInfoErrors, navigate]);

  useEffect(() => {
    const titleFn = () => {
      let base = `${APP_NAME} - `;
      if (isMainInfoFetching) {
        return base + "Загрузка...";
      }
      if (cup && competition && group) {
        const cupPart = cup.title;
        const competiitonPart = getCompetitionStageDescription(competition.stage);
        const groupPart = `${getBowClassDescription(group.bow)} - ${getIdentityDescription(group.identity)}`;
        base += `${cupPart} | ${competiitonPart} | ${groupPart}`;
        switch (groupSection) {
          case "competitors":
            base += " | Участники";
            break;
          case "qualification":
            base += " | Квалификация";
            break;
          case "final":
            base += " | Финал";
            break;
        }
        return base;
      }
      return base + "Ошибка";
    };
    setWebTitle(titleFn());
  }, [competition, cup, group, groupSection, isMainInfoFetching]);

  const groupTitle = group
    ? `${getBowClassDescription(group.bow)} - ${getIdentityDescription(group.identity)}`
    : "Загрузка...";
  const groupSubtitle =
    competition && cup
      ? `${getCompetitionStageDescription(competition.stage)} - ${cup.title}`
      : "Загрузка...";

  return (
    <PageLoader loading={isMainInfoLoading} error={isMainInfoError}>
      <Tabs
        keepMounted={false}
        variant="pills"
        value={groupSection}
        onChange={(section) => section && navigate(section)}
        display="flex"
        flex={1}
        style={{ overflow: "hidden" }}
      >
        <Stack gap="md" display="flex" flex={1} style={{ overflow: "hidden" }}>
          <TabsBar
            title={groupTitle}
            subtitle={groupSubtitle}
            loading={isMainInfoFetching}
            backTo={".."}
            onRefresh={groupSections.controls.onRefresh}
            onExport={groupSections.controls.onExport}
            onComplete={groupSections.controls.onComplete}
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
          </TabsBar>
          <GroupSectionsContext.Provider
            value={{ context: groupSections, setContext: setGroupSections }}
          >
            <Tabs.Panel value="competitors">
              <CompetitorsSection groupId={groupId} />
            </Tabs.Panel>
            <Tabs.Panel value="qualification">
              <QualificationSection groupId={groupId} />
            </Tabs.Panel>
            <Tabs.Panel value="final" display="flex" flex={1} style={{ overflow: "hidden" }}>
              <Text>{"PLACEHOLDER"}</Text>
              {/* <FinalPanel groupInfo={info} setGroupControls={setControls} />*/}
            </Tabs.Panel>
          </GroupSectionsContext.Provider>
        </Stack>
      </Tabs>
    </PageLoader>
  );
}
