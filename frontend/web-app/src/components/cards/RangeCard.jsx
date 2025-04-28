import { Fragment, useState } from "react";
import { IconCheck, IconCircleDashedCheck, IconEdit, IconX } from "@tabler/icons-react";
import {
  ActionIcon,
  Card,
  Divider,
  Group,
  TextInput,
  Tooltip,
  useMantineTheme,
} from "@mantine/core";
import useShotsForm from "../../hooks/useShotsForm";
import BasicRangeCard from "./BasicRangeCard";
import ShotCard from "./ShotCard";

export default function RangeCard({
  range,
  rangeSize,
  rangeRegex,
  loading,
  onShotsSubmit,
  onComplete,
}) {
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

  const renderShots = [...shots]
    .sort((a, b) => a.shotOrdinal - b.shotOrdinal)
    .map((shot, index) => (
      <Fragment key={range.id + shot.shotOrdinal}>
        {index !== 0 && <ShotsDivider />}
        <ShotsFormCard shot={shot} shotsForm={shotsForm} index={index} editing={isEditing} />
      </Fragment>
    ));

  const actionsOnSubmit = async (shotsFormValues) => {
    if (await onShotsSubmit(shotsFormValues.shots)) {
      setIsEditing(false);
    }
  };

  const activeState = range.isActive ? false : range.shots == null ? undefined : true;

  return (
    <form onSubmit={shotsForm.onSubmit(actionsOnSubmit)}>
      <BasicRangeCard
        active={activeState}
        title={"Серия " + range.rangeOrdinal + ":"}
        loading={loading}
      >
        <Group gap="sm" wrap="nowrap">
          {renderShots}
        </Group>
        <ShotsFormControls
          editing={isEditing}
          disabledComplete={!range.isActive}
          onEdit={() => {
            shotsForm.setValues({ shots: [...shots] });
            setIsEditing(true);
          }}
          onComplete={onComplete}
          onCancel={() => {
            setIsEditing(false);
            shotsForm.reset();
          }}
        />
      </BasicRangeCard>
    </form>
  );
}

function ShotsDivider() {
  return <Divider orientation="vertical" size="md" my="xs" />;
}

function ShotsFormCard({ shot, shotsForm, editing }) {
  const theme = useMantineTheme();
  const index = shot.shotOrdinal - 1;

  return (
    <ShotCard score={shot.score} editing={editing}>
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
    </ShotCard>
  );
}

function ShotsFormControls({ editing, disabledComplete, onEdit, onComplete, onCancel }) {
  return (
    <Card p="xs">
      <Group gap="sm">
        {!editing ? (
          <>
            <Tooltip label="Редактировать">
              <ActionIcon onClick={onEdit}>
                <IconEdit />
              </ActionIcon>
            </Tooltip>
            <Tooltip label="Завершить">
              <ActionIcon disabled={disabledComplete} onClick={onComplete}>
                <IconCircleDashedCheck />
              </ActionIcon>
            </Tooltip>
          </>
        ) : (
          <>
            <ActionIcon onClick={onCancel}>
              <IconX />
            </ActionIcon>
            <ActionIcon type="submit">
              <IconCheck />
            </ActionIcon>
          </>
        )}
      </Group>
    </Card>
  );
}
