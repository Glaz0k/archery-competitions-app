import { useEffect, useState } from "react";
import { IconPlus, IconRefresh } from "@tabler/icons-react";
import {
  ActionIcon,
  CloseButton,
  Flex,
  Group,
  Pagination,
  rem,
  Skeleton,
  Stack,
  Text,
  TextInput,
  Title,
  useMantineTheme,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { deleteCup, getCups, postCup } from "../../api/cups";
import { LinkCard, LinkCardSkeleton } from "../cards/LinkCard";
import EmptyCardSpace from "../misc/EmptyCardSpace";
import { CupAddModal, CupDeleteModal } from "../modals/CupModals";

export default function CupsPage() {
  const [cups, setCups] = useState([]);
  const [cupDeletionId, setCupDeletionId] = useState(null);

  const [cupsLoading, setCupsLoading] = useState(true);
  const [cupSubmitting, setCupSubmitting] = useState(false);
  const [cupDeletion, setCupDeletion] = useState(false);

  const [isOpenedAdd, addControl] = useDisclosure(false);
  const [isOpenedDel, delControl] = useDisclosure(false);

  const [searchTerm, setSearchTerm] = useState("");
  const [activePage, setActivePage] = useState(1);

  const theme = useMantineTheme();

  const readCups = async () => {
    try {
      setCupsLoading(true);
      const data = await getCups();
      setCups(data);
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      setCupsLoading(false);
    }
    return true;
  };

  const createCup = async (newCup) => {
    try {
      setCupSubmitting(true);
      const createdCup = await postCup(newCup);
      setCups([createdCup, ...cups]);
      addControl.close();
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      setCupSubmitting(false);
    }
    return true;
  };

  const removeCup = async () => {
    try {
      setCupDeletion(true);
      await deleteCup(cupDeletionId);
      setCups(cups.filter((cup) => cup.id !== cupDeletionId));
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      setCupDeletion(false);
    }
    return true;
  };

  const handleRefresh = async () => {
    if (await readCups()) {
      setActivePage(1);
    }
  };

  const handleSearch = (term) => {
    setSearchTerm(term);
    setActivePage(1);
  };

  const handleCupSubmition = async (cupFormValues) => {
    if (await createCup(cupFormValues)) {
      setActivePage(1);
      return true;
    }
  };

  const confirmCupDeletion = (cupId) => {
    setCupDeletionId(cupId);
    delControl.open();
  };

  const denyCupDeletion = () => {
    setCupDeletionId(null);
  };

  const handleCupDeletion = async () => {
    if (await removeCup()) {
      setActivePage(1);
      delControl.close();
    }
  };

  useEffect(() => {
    readCups();
  }, []);

  const cupsPerPage = 5;

  const filteredCups = cups.filter((cup) =>
    cup.title.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const paginatedCups = filteredCups.slice(
    (activePage - 1) * cupsPerPage,
    activePage * cupsPerPage
  );

  const totalPages = Math.ceil(filteredCups.length / cupsPerPage);

  return (
    <>
      <CupAddModal
        opened={isOpenedAdd}
        onClose={addControl.close}
        handleSubmit={handleCupSubmition}
        loading={cupSubmitting}
      />
      <CupDeleteModal
        opened={isOpenedDel}
        onClose={delControl.close}
        onDeny={denyCupDeletion}
        onConfirm={handleCupDeletion}
        loading={cupDeletion}
      />

      <Flex direction="column" flex={1}>
        <Group>
          <Title order={2} flex={1}>
            Кубки
          </Title>
          <TextInput
            placeholder="Название"
            value={searchTerm}
            onChange={(e) => handleSearch(e.currentTarget.value)}
            rightSection={<CloseButton onClick={() => handleSearch("")} />}
          />
          <ActionIcon onClick={handleRefresh}>
            <IconRefresh />
          </ActionIcon>
          <ActionIcon onClick={addControl.open}>
            <IconPlus />
          </ActionIcon>
        </Group>
        <Stack flex={1}>
          {cupsLoading ? (
            Array(cupsPerPage)
              .fill(0)
              .map((_, index) => (
                <LinkCardSkeleton key={index} isDelete>
                  <Skeleton height={rem(theme.fontSizes.md)} width={400} />
                  <Skeleton height={rem(theme.fontSizes.md)} width={150} />
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
                <Text>Адрес: {cup.address != null ? cup.address : "Не указан"}</Text>
                <Text>Сезон: {cup.season != null ? cup.season : "Не указан"}</Text>
              </LinkCard>
            ))
          ) : (
            <EmptyCardSpace label="Кубки не найдены" />
          )}
        </Stack>
        <Pagination
          value={activePage}
          onChange={setActivePage}
          total={cupsLoading ? 0 : totalPages}
        />
      </Flex>
    </>
  );
}
