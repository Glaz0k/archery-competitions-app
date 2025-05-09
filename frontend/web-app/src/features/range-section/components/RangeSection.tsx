import { useState } from "react";
import { Fragment } from "react/jsx-runtime";
import { Divider, Group, TextInput, useMantineTheme } from "@mantine/core";
import { RangeType, type Shot } from "../../../entities";
import { NO_SCORE_VALUE } from "../../constants";
import { useShotsForm } from "../hooks/useShotsForm";
import { ScoreCard } from "../widgets/ScoreCard";
import { ShotsFormControls } from "./ShotsFormControls";

export interface RangeSectionProps {
  shots: Shot[];
  type: RangeType;
  editFn: (shots: Shot[]) => unknown;
  completeFn: () => unknown;
  active: boolean;
}

export function RangeSection({ shots, type, editFn, completeFn, active }: RangeSectionProps) {
  const shotRg = type === RangeType.ONE_TEN ? /^(M|[1-9]|10|X)$/ : /^(M|[6-9]|10|X)$/;
  const shotsForm = useShotsForm(shots, shotRg);

  const [isEditing, setEditing] = useState<boolean>(false);

  const renderShots = shots.map((shot, index) => (
    <Fragment key={`shot${shot.ordinal}$score${shot.score || NO_SCORE_VALUE}`}>
      {index !== 0 && <ShotsDivider />}
      <ShotsFormCard shot={shot} index={index} form={shotsForm} editing={isEditing} />
    </Fragment>
  ));

  return (
    <form
      onSubmit={shotsForm.onSubmit((shots) => {
        editFn(shots);
      })}
    >
      <Group gap="sm" wrap="nowrap">
        {renderShots}
      </Group>
      <ShotsFormControls
        editing={isEditing}
        active={active}
        onEdit={() => {
          shotsForm.setValues({
            shots,
          });
          setEditing(true);
        }}
        onCancelEdit={() => {
          setEditing(false);
          shotsForm.reset();
        }}
        onComplete={completeFn}
      />
    </form>
  );
}

function ShotsDivider() {
  return <Divider orientation="vertical" size="md" my="xs" />;
}

interface ShotsFormCardProps {
  shot: Shot;
  index: number;
  form: ReturnType<typeof useShotsForm>;
  editing: boolean;
}

function ShotsFormCard({ shot, index, form, editing }: ShotsFormCardProps) {
  const theme = useMantineTheme();
  const formIndex = `shots.${index}.score`;
  return (
    <ScoreCard score={shot.score} editing={editing}>
      <TextInput
        p={0}
        variant="unstyled"
        styles={{
          input: {
            font: theme.headings.fontFamily,
            fontSize: theme.headings.sizes.h2.fontSize,
            fontWeight: theme.headings.sizes.h2.fontWeight,
            color: form.errors[formIndex] === undefined ? theme.white : theme.colors.red[6],
          },
        }}
        key={form.key(formIndex)}
        {...form.getInputProps(formIndex, { withError: false })}
      />
    </ScoreCard>
  );
}
