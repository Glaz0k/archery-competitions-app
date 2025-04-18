import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  Center,
  CloseButton,
  Flex,
  Pagination,
  rem,
  Skeleton,
  Stack,
  Text,
  TextInput,
  useMantineTheme,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { deleteCup, getCups, postCup } from "../../api/cups";
import { CUP_QUERY_KEYS } from "../../api/queryKeys";
import MainBar from "../bars/MainBar";
import { LinkCard, LinkCardSkeleton } from "../cards/LinkCard";
import NotFoundCard from "../cards/NotFoundCard";
import AddCupModal from "../modals/cup/AddCupModal";
import DeleteCupModal from "../modals/cup/DeleteCupModal";

const CUPS_PER_PAGE = 5;

export default function CupsPage() {
  const queryClient = useQueryClient();
  const theme = useMantineTheme();

  const [cupDeletionId, setCupDeletionId] = useState(null);

  const [isOpenedAdd, addControl] = useDisclosure(false);
  const [isOpenedDel, delControl] = useDisclosure(false);

  const [searchTerm, setSearchTerm] = useState("");
  const [activePage, setActivePage] = useState(1);

  const { mutateAsync: createCup, isPending: isCupSubmitting } = useMutation({
    mutationFn: (newCup) => postCup(newCup),
    onSuccess: (createdCup) => {
      queryClient.setQueryData(CUP_QUERY_KEYS.all, (old) => [createdCup, ...(old || [])]);
      addControl.close();
      setActivePage(1);
    },
  });

  const {
    data: cups,
    isFetching: isCupsLoading,
    refetch: refetchCups,
  } = useQuery({
    queryKey: CUP_QUERY_KEYS.all,
    queryFn: getCups,
    initialData: [],
  });

  const { mutate: removeCup, isPending: isCupDeleting } = useMutation({
    mutationFn: () => deleteCup(cupDeletionId),
    onSuccess: () => {
      queryClient.setQueryData(CUP_QUERY_KEYS.all, (old) =>
        old.filter((cup) => cup.id !== cupDeletionId)
      );
      denyCupDeletion();
      setActivePage(1);
    },
  });

  const handleSearch = (term) => {
    setSearchTerm(term);
    setActivePage(1);
  };

  const handleRefresh = () => {
    refetchCups().then(() => setActivePage(1));
  };

  const confirmCupDeletion = (cupId) => {
    setCupDeletionId(cupId);
    delControl.open();
  };

  const denyCupDeletion = () => {
    delControl.close();
    setCupDeletionId(null);
  };

  const filteredCups = cups.filter((cup) =>
    cup.title.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const paginatedCups = filteredCups.slice(
    (activePage - 1) * CUPS_PER_PAGE,
    activePage * CUPS_PER_PAGE
  );

  const totalPages = Math.ceil(filteredCups.length / CUPS_PER_PAGE);

  return (
    <>
      <AddCupModal
        isOpened={isOpenedAdd}
        onClose={addControl.close}
        onSubmit={createCup}
        isLoading={isCupSubmitting}
      />
      <DeleteCupModal
        isOpened={isOpenedDel}
        onClose={denyCupDeletion}
        onConfirm={removeCup}
        isLoading={isCupDeleting}
      />

      <Flex direction="column" flex={1} gap="lg">
        <MainBar title={"Кубки"} onRefresh={handleRefresh} onAdd={addControl.open}>
          <TextInput
            placeholder="Название"
            value={searchTerm}
            onChange={(e) => handleSearch(e.currentTarget.value)}
            rightSection={<CloseButton onClick={() => handleSearch("")} />}
          />
        </MainBar>
        <Stack flex={1} gap="md">
          {isCupsLoading ? (
            Array(CUPS_PER_PAGE)
              .fill(0)
              .map((_, index) => (
                <LinkCardSkeleton key={index} isDelete>
                  <Skeleton height={rem(theme.fontSizes.sm)} width={400} />
                  <Skeleton height={rem(theme.fontSizes.sm)} width={150} />
                </LinkCardSkeleton>
              ))
          ) : paginatedCups.length > 0 ? (
            paginatedCups.map((cup, index) => (
              <LinkCard
                key={index}
                title={cup.title}
                onDelete={() => confirmCupDeletion(cup.id)}
                to={"/cups/" + cup.id}
              >
                <Text size="sm">
                  {"Адрес: "}
                  {cup.address !== null ? cup.address : "Не указан"}
                </Text>
                <Text size="sm">
                  {"Сезон: "}
                  {cup.season !== null ? cup.season : "Не указан"}
                </Text>
              </LinkCard>
            ))
          ) : (
            <NotFoundCard label="Кубки не найдены" />
          )}
        </Stack>
        <Center>
          <Pagination
            value={activePage}
            onChange={setActivePage}
            total={isCupsLoading ? 0 : totalPages}
          />
        </Center>
      </Flex>
    </>
  );
}
