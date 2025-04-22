import { Fragment, useState } from "react";
import {
  IconCheck,
  IconCircleDashedCheck,
  IconCircleFilled,
  IconEdit,
  IconX,
} from "@tabler/icons-react";
import {
  ActionIcon,
  AspectRatio,
  Card,
  Group,
  LoadingOverlay,
  TextInput,
  ThemeIcon,
  Title,
  useMantineTheme,
} from "@mantine/core";
import useShotsForm from "../../hooks/useShotsForm";
import { ShotCard, ShotDivider } from "./ShotCard";

export default function RangeCard({
  range,
  rangeSize,
  rangeRegex,
  loading,
  onShotsSubmit,
  onComplete,
}) {
  const theme = useMantineTheme();
  const [isEditing, setIsEditing] = useState(false);
  const shots =
    range.shots != null
      ? [...range.shots]
      : Array(rangeSize)
          .fill(0)
          .map((_, index) => {
            return {
              shotOrdinal: index + 1,
              score: null,
            };
          });

  const shotsForm = useShotsForm({ initialShots: [...shots], shotRegex: rangeRegex });

  const shotCards = [...shots]
    .sort((a, b) => a.shotOrdinal - b.shotOrdinal)
    .map((shot, index) => (
      <Fragment key={shot.shotOrdinal}>
        {index !== 0 && <ShotDivider />}
        <AspectRatio ratio={1}>
          <ShotCard>
            {isEditing ? (
              <TextInput
                p={0}
                variant="unstyled"
                styles={{
                  input: {
                    font: theme.headings.fontFamily,
                    fontSize: theme.headings.sizes.h2.fontSize,
                    fontWeight: theme.headings.sizes.h2.fontWeight,
                    color:
                      shotsForm.errors["shots." + index + ".score"] == undefined
                        ? theme.white
                        : theme.colors.red[6],
                  },
                }}
                key={shotsForm.key("shots." + index + ".score")}
                {...shotsForm.getInputProps("shots." + index + ".score", { withError: false })}
              />
            ) : (
              <Title order={2}>{shot.score || "-"}</Title>
            )}
          </ShotCard>
        </AspectRatio>
      </Fragment>
    ));

  const actionsOnSubmit = async (addFormValues) => {
    if (await onShotsSubmit(addFormValues)) {
      setIsEditing(false);
    }
  };

  return (
    <Card bg="gray.0" bd="5px solid #E0E0E0" c={theme.black}>
      <LoadingOverlay visible={loading} />
      <form onSubmit={shotsForm.onSubmit(actionsOnSubmit)}>
        <Group gap="md">
          <ThemeIcon
            color={
              range.isActive
                ? theme.colors.yellow[6]
                : range.shots == null
                  ? theme.colors.gray[6]
                  : theme.colors.green[6]
            }
          >
            <IconCircleFilled />
          </ThemeIcon>
          <Title order={2} flex={1}>
            {"Серия " + range.rangeOrdinal + ":"}
          </Title>
          <Group gap="sm" wrap="nowrap">
            {shotCards}
          </Group>
          <Card p="xs">
            <Group gap="sm">
              {!isEditing ? (
                <>
                  <ActionIcon
                    onClick={() => {
                      shotsForm.setInitialValues({ shots: [...shots] });
                      setIsEditing(true);
                    }}
                  >
                    <IconEdit />
                  </ActionIcon>
                  <ActionIcon disabled={!range.isActive} onClick={() => onComplete()}>
                    <IconCircleDashedCheck />
                  </ActionIcon>
                </>
              ) : (
                <>
                  <>
                    <ActionIcon
                      onClick={() => {
                        setIsEditing(false);
                        shotsForm.reset();
                      }}
                    >
                      <IconX />
                    </ActionIcon>
                    <ActionIcon type="submit">
                      <IconCheck />
                    </ActionIcon>
                  </>
                </>
              )}
            </Group>
          </Card>
        </Group>
      </form>
    </Card>
  );
}
