import { useState } from "react";
import { useNavigate, useParams } from "react-router";
import { Stack } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import {
  useCreateIndividualGroup,
  useDeleteIndividualGroup,
  useIndividualGroups,
} from "../../../api";
import { GroupState } from "../../../entities";
import {
  AddIndividualGroupModal,
  DeleteIndividualGroupModal,
  GroupStateBadge,
  useGroupFilter,
} from "../../../features";
import { getBowClassDescription, getIdentityDescription } from "../../../utils";
import { CenterCard, EntityCard, EntityCardSkeleton, TopBar } from "../../../widgets";

const SKELETON_LENGTH = 6;

function isExported(state?: GroupState) {
  switch (state) {
    case GroupState.CREATED:
    case GroupState.QUAL_START:
      return false;
    default:
      return true;
  }
}

export default function GroupsContent() {
  const { competitionId: paramCompetitionId } = useParams();
  const competitionId = Number(paramCompetitionId);

  const navigate = useNavigate();

  const [isOpenedGroupAdd, groupAddControl] = useDisclosure(false);
  const [isOpenedGroupDel, groupDelControl] = useDisclosure(false);

  const [groupDeletingId, setGroupDelitingId] = useState<number | null>(null);

  const { filter } = useGroupFilter();

  const {
    data: groups,
    isFetching: isGroupsLoading,
    refetch: refetchGroups,
    isError: isGroupsError,
    error,
  } = useIndividualGroups(competitionId);

  const { mutateAsync: createGroup, isPending: isGroupSubmitting } = useCreateIndividualGroup(
    groupAddControl.close
  );

  const { mutate: deleteGroup, isPending: isGroupDeleting } = useDeleteIndividualGroup(() => {
    groupDelControl.close();
    setGroupDelitingId(null);
  });

  const handleGroupDeleting = (id: number) => {
    setGroupDelitingId(id);
    groupDelControl.open();
  };

  const handleExport = (id: number) => {
    console.warn(`handleExport temporary unavailable: ${id}`);
  };

  const filteredGroups = groups.filter((group) => {
    if (filter.identity && filter.identity !== group.identity) {
      return false;
    }
    if (filter.state && filter.state !== group.state) {
      return false;
    }
    if (filter.bow && filter.bow !== group.bow) {
      return false;
    }
    return true;
  });

  let renderContent;
  if (isGroupsLoading) {
    renderContent = Array(SKELETON_LENGTH)
      .fill(0)
      .map((_, index) => <EntityCardSkeleton key={index} tagged exported deleted />);
  } else if (isGroupsError) {
    console.error(error);
    renderContent = <CenterCard label={"Произошла ошибка"} />;
  } else if (filteredGroups.length === 0) {
    renderContent = <CenterCard label={"Дивизионы не найдены"} />;
  } else {
    renderContent = filteredGroups.map(({ id, bow, identity, state }, index) => {
      return (
        <EntityCard
          key={index}
          title={getBowClassDescription(bow) + " - " + getIdentityDescription(identity)}
          tag={<GroupStateBadge state={state} />}
          to={"individual-groups/" + id}
          onDelete={() => handleGroupDeleting(id)}
          onExport={isExported(state) ? () => handleExport(id) : undefined}
        />
      );
    });
  }

  return (
    <>
      <AddIndividualGroupModal
        opened={isOpenedGroupAdd}
        onClose={groupAddControl.close}
        onSubmit={({ bow, identity }) => createGroup([competitionId, { bow, identity }])}
        loading={isGroupSubmitting}
      />
      <DeleteIndividualGroupModal
        opened={isOpenedGroupDel}
        onClose={groupDelControl.close}
        onConfirm={() => deleteGroup(groupDeletingId!)}
        loading={isGroupDeleting}
      />
      <Stack flex={1} h="100%" gap="lg" miw={500}>
        <TopBar
          title={"Дивизионы"}
          onRefresh={refetchGroups}
          onAdd={groupAddControl.open}
          onBack={() => navigate("..")}
        />
        <Stack flex={1}>{renderContent}</Stack>
      </Stack>
    </>
  );
}
