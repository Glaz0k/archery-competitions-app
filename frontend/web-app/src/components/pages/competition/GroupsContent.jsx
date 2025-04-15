import { useState } from "react";
import { IconRefresh } from "@tabler/icons-react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useOutletContext, useParams } from "react-router";
import { ActionIcon, Stack } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { getIndividualGroups, postIndividualGroup } from "../../../api/competitions";
import { deleteIndividualGroup } from "../../../api/individualGroups";
import { INDIVIDUAL_GROUP_QUERY_KEYS } from "../../../api/queryKeys";
import BowClass from "../../../enums/BowClass";
import GroupGender from "../../../enums/GroupGender";
import GroupState from "../../../enums/GroupState";
import MainBar from "../../bars/MainBar";
import { LinkCard, LinkCardSkeleton } from "../../cards/LinkCard";
import EmptyCardSpace from "../../misc/EmptyCardSpace";
import GroupStateBadge from "../../misc/GroupStateBadge";
import AddIndividualGroupModal from "../../modals/AddIndividualGroupModal";
import DeleteIndividualGroupModal from "../../modals/DeleteIndividualGroupModal";

const SKELETON_LENGTH = 7;

function isExported(groupState) {
  switch (groupState) {
    case GroupState.CREATED:
    case GroupState.QUAL_START:
      return false;
    default:
      return true;
  }
}

export default function GroupsContent() {
  const groupsFilter = useOutletContext();
  const queryClient = useQueryClient();
  const { competitionId } = useParams();

  const [isOpenedGroupAdd, groupAddControl] = useDisclosure(false);
  const [isOpenedGroupDel, groupDelControl] = useDisclosure(false);

  const [groupDeletingId, setGroupDelitingId] = useState(null);

  const { mutateAsync: createGroup, isPending: isGroupSubmitting } = useMutation({
    mutationFn: (newGroup) => postIndividualGroup(competitionId, newGroup),
    onSuccess: (newGroup) => {
      queryClient.setQueryData(
        INDIVIDUAL_GROUP_QUERY_KEYS.allByCompetition(competitionId),
        (old) => [newGroup, ...(old || [])]
      );
      groupAddControl.close();
    },
  });

  const {
    data: groups,
    isFetching: isGroupsLoading,
    refetch: refetchGroups,
  } = useQuery({
    queryKey: INDIVIDUAL_GROUP_QUERY_KEYS.allByCompetition(competitionId),
    queryFn: () => getIndividualGroups(competitionId),
    initialData: [],
  });

  const { mutate: removeGroup, isPending: isGroupDeleting } = useMutation({
    mutationFn: () => deleteIndividualGroup(groupDeletingId),
    onSuccess: () => {
      queryClient.setQueryData(
        INDIVIDUAL_GROUP_QUERY_KEYS.allByCompetition(competitionId),
        (old) => {
          return old.filter((group) => group.id !== groupDeletingId);
        }
      );
      groupDelControl.close();
      setGroupDelitingId(null);
    },
  });

  const handleGroupDeleting = (id) => {
    setGroupDelitingId(id);
    groupDelControl.open();
  };

  const handleExport = (_id) => {
    console.warn("handleExport temporary unavailable");
  };

  const filteredGroups = groups.filter((group) => {
    if (groupsFilter.bow !== "default" && group.bow !== BowClass.valueOf(groupsFilter.bow)) {
      return false;
    }
    if (
      groupsFilter.identity !== "default" &&
      group.identity !== GroupGender.valueOf(groupsFilter.identity)
    ) {
      return false;
    }
    if (
      groupsFilter.state !== "default" &&
      group.state !== GroupState.valueOf(groupsFilter.state)
    ) {
      return false;
    }
    return true;
  });

  return (
    <>
      <AddIndividualGroupModal
        isOpened={isOpenedGroupAdd}
        onClose={groupAddControl.close}
        onSubmit={createGroup}
        isLoading={isGroupSubmitting}
      />
      <DeleteIndividualGroupModal
        isOpened={isOpenedGroupDel}
        onClose={groupDelControl.close}
        onConfirm={removeGroup}
        isLoading={isGroupDeleting}
      />
      <Stack flex={1} h="100%">
        <MainBar
          title={"Индивидуальные группы"}
          onRefresh={refetchGroups}
          onAdd={groupAddControl.open}
        ></MainBar>
        <Stack flex={1}>
          {isGroupsLoading ? (
            Array(SKELETON_LENGTH)
              .fill(0)
              .map((_, index) => <LinkCardSkeleton key={index} isTagged isExport isDelete />)
          ) : filteredGroups.length !== 0 ? (
            filteredGroups.map(({ id, bow, identity, state }, index) => {
              return (
                <LinkCard
                  key={index}
                  title={bow.textValue + " - " + identity.textValue}
                  tag={<GroupStateBadge state={state} />}
                  to={"individual-groups/" + id}
                  onDelete={() => handleGroupDeleting(id)}
                  onExport={isExported(state) ? () => handleExport(id) : null}
                />
              );
            })
          ) : (
            <EmptyCardSpace label={"Группы не найдены"} />
          )}
        </Stack>
      </Stack>
    </>
  );
}
