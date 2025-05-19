import { useState } from "react";
import { Group, LoadingOverlay, Stack, Tabs, Title } from "@mantine/core";
import {
  useCompleteSectionRange,
  useEditSectionRange,
  useQualificationSection,
  useSectionRangeGroup,
  useSectionRound,
} from "../../../api";
import { RangeType, type Range } from "../../../entities";
import { ControlsCard, PageLoader, SideBar, TopBar } from "../../../widgets";
import { RangeCard, RangeSection } from "../../range-section";

export interface SectionTabProps {
  sectionId: number;
}

export function SectionTab({ sectionId }: SectionTabProps) {
  const [selectedRoundOrd, setSelectedRoundOrd] = useState<number | null>(null);

  const {
    data: section,
    isFetching: isSectionFetching,
    isLoading: isSectionLoading,
    isError: isSectionError,
  } = useQualificationSection(sectionId);

  return (
    <PageLoader loading={isSectionLoading} error={isSectionError}>
      {section && (
        <Group flex={1} align="start">
          <SideBar>
            <ControlsCard pos="relative">
              <LoadingOverlay visible={isSectionFetching} />
              <Stack>
                <Title order={2}>{section.competitor.fullName}</Title>
                <Tabs
                  orientation="vertical"
                  variant="pills"
                  value={selectedRoundOrd ? String(selectedRoundOrd) : null}
                  onChange={(value) => setSelectedRoundOrd(value ? Number(value) : null)}
                >
                  <Tabs.List>
                    {section.rounds
                      .map((round) => round.ordinal)
                      .sort()
                      .map((ordinal) => (
                        <Tabs.Tab key={`${section.id}$round${ordinal}`} value={String(ordinal)}>
                          <Title order={3}>{`Раунд ${ordinal}`}</Title>
                        </Tabs.Tab>
                      ))}
                  </Tabs.List>
                </Tabs>
              </Stack>
            </ControlsCard>
          </SideBar>
          {selectedRoundOrd && <RoundTab sectionId={sectionId} roundOrd={selectedRoundOrd} />}
        </Group>
      )}
    </PageLoader>
  );
}

interface RoundTabProps {
  sectionId: number;
  roundOrd: number;
}

function RoundTab({ sectionId, roundOrd }: RoundTabProps) {
  const {
    data: round,
    isLoading: isRoundLoading,
    isError: isRoundError,
    isFetching: isRoundFetching,
    refetch: refetchRound,
  } = useSectionRound(sectionId, roundOrd);

  const {
    data: rangeGroup,
    isLoading: isRangeGroupLoading,
    isError: isRangeGroupError,
    isFetching: isRangeGroupFetching,
    refetch: refetchRangeGroup,
  } = useSectionRangeGroup(sectionId, roundOrd);

  return (
    <PageLoader
      loading={isRoundLoading || isRangeGroupLoading}
      error={isRoundError || isRangeGroupError}
    >
      {round && rangeGroup && (
        <Stack flex={1}>
          <TopBar
            title={`Раунд ${roundOrd}${round.isActive ? " | Активен" : ""}`}
            onRefresh={() => {
              refetchRound();
              refetchRangeGroup();
            }}
            loading={isRoundFetching || isRangeGroupFetching}
          />
          {[...rangeGroup.ranges]
            .sort(({ ordinal: a }, { ordinal: b }) => a - b)
            .map((range) => (
              <SectionRangeCard
                key={`${sectionId}$round${roundOrd}$range${range.ordinal}`}
                sectionId={sectionId}
                roundOrd={round.ordinal}
                range={range}
                rangeSize={rangeGroup.rangeSize}
                rangeType={rangeGroup.type}
              />
            ))}
        </Stack>
      )}
    </PageLoader>
  );
}

interface SectionRangeCardProps {
  sectionId: number;
  roundOrd: number;
  range: Range;
  rangeSize: number;
  rangeType: RangeType;
}

function SectionRangeCard({
  sectionId,
  roundOrd,
  range,
  rangeSize,
  rangeType,
}: SectionRangeCardProps) {
  const { mutate: editRange, isPending: isRangeEditing } = useEditSectionRange();
  const { mutate: completeRange, isPending: isRangeCompleting } = useCompleteSectionRange();

  const shotOrdinals = [...Array(rangeSize).keys()].map((i) => i + 1);
  const shots = shotOrdinals.map((ord) => {
    let shot = range.shots?.find((shot) => shot.ordinal === ord);
    if (!shot) {
      shot = {
        ordinal: ord,
        score: null,
      };
    }
    return shot;
  });

  return (
    <RangeCard
      active={range.isActive}
      title={`Серия ${range.ordinal}`}
      loading={isRangeEditing || isRangeCompleting}
    >
      <RangeSection
        shots={shots}
        type={rangeType}
        editFn={(editedShots) =>
          editRange([sectionId, roundOrd, { ordinal: range.ordinal, shots: editedShots }])
        }
        completeFn={() => completeRange([sectionId, roundOrd, range.ordinal])}
        active={range.isActive}
      />
    </RangeCard>
  );
}
