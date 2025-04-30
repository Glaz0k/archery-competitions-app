import { useState } from "react";
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
import { useCreateCup, useCups, useDeleteCup } from "../../api";
import MainBar from "../bars/MainBar";
import { LinkCard, LinkCardSkeleton } from "../cards/LinkCard";
import NotFoundCard from "../cards/NotFoundCard";
import AddCupModal from "../modals/cup/AddCupModal";
import DeleteCupModal from "../modals/cup/DeleteCupModal";

const CUPS_PER_PAGE = 3;

export default function CupsPage() {
  const theme = useMantineTheme();

  const [cupDeletionId, setCupDeletionId] = useState(null);

  const [isOpenedAdd, addControl] = useDisclosure(false);
  const [isOpenedDel, delControl] = useDisclosure(false);

  const [searchTerm, setSearchTerm] = useState("");
  const [activePage, setActivePage] = useState(1);

  const { mutateAsync: createCup, isPending: isCupSubmitting } = useCreateCup(() => {
    addControl.close();
    setActivePage(1);
  });

  const {
    data: cups,
    isFetching: isCupsLoading,
    refetch: refetchCups,
    isError: isCupsError,
    error: cupsError,
  } = useCups();

  const { mutate: removeCup, isPending: isCupDeleting } = useDeleteCup(() => {
    denyCupDeletion();
    setActivePage(1);
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

  let renderContent;
  if (isCupsLoading) {
    renderContent = Array(CUPS_PER_PAGE)
      .fill(0)
      .map((_, index) => (
        <LinkCardSkeleton key={index} isDelete>
          <Skeleton height={rem(theme.fontSizes.sm)} width={400} />
          <Skeleton height={rem(theme.fontSizes.sm)} width={150} />
        </LinkCardSkeleton>
      ));
  } else if (isCupsError) {
    console.error(cupsError.name + "\n" + cupsError.message);
    renderContent = <NotFoundCard label="Произошла ошибка" />;
  } else if (paginatedCups.length === 0) {
    renderContent = <NotFoundCard label="Кубки не найдены" />;
  } else {
    renderContent = paginatedCups.map((cup, index) => (
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
    ));
  }

  return (
    <>
      <AddCupModal
        opened={isOpenedAdd}
        onClose={addControl.close}
        onSubmit={async (values) => {
          console.log(values);
          await createCup(values);
        }}
        loading={isCupSubmitting}
      />
      <DeleteCupModal
        isOpened={isOpenedDel}
        onClose={denyCupDeletion}
        onConfirm={() => removeCup(Number(cupDeletionId))}
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
          {renderContent}
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
