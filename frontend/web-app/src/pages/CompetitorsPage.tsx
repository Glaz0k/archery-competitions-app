import { forwardRef, useEffect, useState } from "react";
import { IconEdit, IconLockCog, IconTrashX, IconUser } from "@tabler/icons-react";
import { format } from "date-fns";
import {
  ActionIcon,
  Card,
  Center,
  CloseButton,
  createPolymorphicComponent,
  Flex,
  Group,
  HoverCard,
  Loader,
  Pagination,
  Stack,
  Text,
  TextInput,
  ThemeIcon,
  Tooltip,
  useMantineTheme,
  type CardProps,
} from "@mantine/core";
import { useDocumentTitle } from "@mantine/hooks";
import { useCompetitors } from "../api";
import { APP_NAME } from "../constants";
import type { Competitor } from "../entities";
import { getBowClassDescription, getGenderDescription, getSportsRankDescription } from "../utils";
import { CenterCard, TopBar } from "../widgets";

const COMPETITORS_PER_PAGE = 10;

export default function CompetitorsPage() {
  const [searchTerm, setSearchTerm] = useState<string>("");
  const [activePage, setActivePage] = useState<number>(1);

  const {
    data: competitors,
    refetch: refetchCompetitors,
    isFetching: isCompetitorsFetching,
    isError: isCompetitorsError,
  } = useCompetitors();

  const handleSearch = (term: string) => {
    setSearchTerm(term);
    setActivePage(1);
  };

  const filteredCompetitors = competitors.filter((competitor) =>
    competitor.fullName.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const paginatedCompetitors = filteredCompetitors.slice(
    (activePage - 1) * COMPETITORS_PER_PAGE,
    activePage * COMPETITORS_PER_PAGE
  );

  const totalPages = Math.ceil(filteredCompetitors.length / COMPETITORS_PER_PAGE);

  let renderContent;
  if (isCompetitorsFetching) {
    renderContent = (
      <Center flex={1}>
        <Loader size="lg" />
      </Center>
    );
  } else if (isCompetitorsError) {
    renderContent = <CenterCard label="Произошла ошибка" />;
  } else if (paginatedCompetitors.length === 0) {
    renderContent = <CenterCard label="Ничего не найдено" />;
  } else {
    renderContent = (
      <>
        <Stack flex={1} gap="md">
          {paginatedCompetitors.map((val) => (
            <CompetitorHoverCard key={val.id} competitor={val} />
          ))}
        </Stack>
        <Center>
          <Pagination value={activePage} onChange={setActivePage} total={totalPages} />
        </Center>
      </>
    );
  }

  const [webTitle, setWebTitle] = useState<string>("");
  useDocumentTitle(webTitle);
  useEffect(() => {
    setWebTitle(`${APP_NAME} | Пользователи`);
  }, []);

  return (
    <Flex direction="row" flex={1} gap="lg">
      <Stack flex={1} gap="lg">
        <TopBar title="Пользователи" onRefresh={refetchCompetitors}>
          <TextInput
            placeholder="Фамилия Имя"
            value={searchTerm}
            onChange={(e) => handleSearch(e.currentTarget.value)}
            rightSection={<CloseButton onClick={() => handleSearch("")} />}
          />
        </TopBar>
        {renderContent}
      </Stack>
    </Flex>
  );
}

function CompetitorHoverCard({ competitor }: { competitor: Competitor }) {
  const theme = useMantineTheme();
  return (
    <HoverCard
      offset={{ mainAxis: 10, crossAxis: 60 }}
      position="bottom-start"
      withArrow
      arrowSize={10}
      arrowOffset={25}
      styles={{
        dropdown: {
          borderColor: theme.colors.secondary![9],
          borderWidth: 2,
        },
        arrow: {
          borderTopWidth: 2,
          borderTopColor: `${theme.colors.secondary![9]}`,
          borderBottomWidth: 2,
          borderBottomColor: `${theme.colors.secondary![9]}`,
          borderLeftWidth: 2,
          borderLeftColor: `${theme.colors.secondary![9]}`,
          borderRightWidth: 2,
          borderRightColor: `${theme.colors.secondary![9]}`,
        },
      }}
    >
      <HoverCard.Target>
        <CompetitorCard competitor={competitor} />
      </HoverCard.Target>
      <HoverCard.Dropdown>
        <Text>{`ID: ${competitor.id}`}</Text>
        <Text>{`Класс лука: ${competitor.bow ? getBowClassDescription(competitor.bow) : "Не указан"}`}</Text>
        <Text>{`Разряд: ${competitor.rank ? getSportsRankDescription(competitor.rank) : "б/р"}`}</Text>
        <Text>{`Регион: ${competitor.region || "Не указан"}`}</Text>
        <Text>{`Федерация: ${competitor.federation || "Не указана"}`}</Text>
        <Text>{`Клуб: ${competitor.club || "Не указан"}`}</Text>
      </HoverCard.Dropdown>
    </HoverCard>
  );
}

interface CompetitorCardProps extends CardProps {
  competitor: Competitor;
}

const CompetitorCard = createPolymorphicComponent<"div", CompetitorCardProps>(
  forwardRef<HTMLDivElement, CompetitorCardProps>(({ competitor, ...others }, ref) => (
    <Card component="div" padding="md" ref={ref} {...others}>
      <Group gap="md">
        <ThemeIcon>
          <IconUser />
        </ThemeIcon>
        <Stack flex={1} gap={0}>
          <Text>{competitor.fullName}</Text>
          <Text size="sm">
            {`${getGenderDescription(competitor.identity)}, ${format(competitor.birthDate, "dd.MM.yyyy")}`}
          </Text>
        </Stack>
        <Tooltip label="Редактировать профиль">
          <ActionIcon>
            <IconEdit />
          </ActionIcon>
        </Tooltip>
        <Tooltip label="Восстановить пароль">
          <ActionIcon>
            <IconLockCog />
          </ActionIcon>
        </Tooltip>
        <Tooltip label="Удалить учётную запись">
          <ActionIcon>
            <IconTrashX />
          </ActionIcon>
        </Tooltip>
      </Group>
    </Card>
  ))
);
