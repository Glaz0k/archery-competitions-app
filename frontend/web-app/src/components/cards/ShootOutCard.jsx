import { useEffect, useState } from "react";
import { IconCheck, IconEdit, IconX } from "@tabler/icons-react";
import { ActionIcon, Card, Group, TextInput, Tooltip, useMantineTheme } from "@mantine/core";
import { matches, useForm } from "@mantine/form";
import BasicRangeCard from "./BasicRangeCard";
import ShotCard from "./ShotCard";

export default function ShootOutCard({ shootOut, shootOutRegex, loading, onShootOutSubmit }) {
  const theme = useMantineTheme();
  const [isEditing, setIsEditing] = useState(false);

  const shootOutForm = useForm({
    mode: "uncontrolled",
    initialValues: {
      fullScore: shootOut.score,
    },
    validateInputOnChange: ["fullScore"],
    validate: {
      fullScore: matches(shootOutRegex, "Неверный счёт"),
    },
  });

  const actionsOnSubmit = async (shotsFormValues) => {
    if (await onShootOutSubmit(scoreToShootOut(shootOut.id, shotsFormValues.fullScore))) {
      setIsEditing(false);
    }
  };

  useEffect(() => {});

  return (
    <form onSubmit={shootOutForm.onSubmit(actionsOnSubmit)}>
      <BasicRangeCard active={!!shootOut.score} title={"Перестрелка:"} loading={loading}>
        <ShotCard score={shootOutToScore(shootOut)} editing={isEditing}>
          <TextInput
            p={0}
            variant="unstyled"
            styles={{
              input: {
                font: theme.headings.fontFamily,
                fontSize: theme.headings.sizes.h2.fontSize,
                fontWeight: theme.headings.sizes.h2.fontWeight,
                color:
                  shootOutForm.errors["fullScore"] == undefined ? theme.white : theme.colors.red[6],
              },
            }}
            key={shootOutForm.key("fullScore")}
            {...shootOutForm.getInputProps("fullScore", { withError: false })}
          />
        </ShotCard>
        <ShootOutFormControls
          editing={isEditing}
          onEdit={() => {
            shootOutForm.setValues({ fullScore: shootOutToScore(shootOut) });
            setIsEditing(true);
          }}
          onCancel={() => {
            setIsEditing(false);
            shootOutForm.reset();
          }}
        />
      </BasicRangeCard>
    </form>
  );
}

function ShootOutFormControls({ editing, onEdit, onCancel }) {
  return (
    <Card p="xs">
      <Group gap="sm">
        {!editing ? (
          <Tooltip label="Редактировать">
            <ActionIcon onClick={onEdit}>
              <IconEdit />
            </ActionIcon>
          </Tooltip>
        ) : (
          <>
            <ActionIcon type="submit">
              <IconCheck />
            </ActionIcon>
            <ActionIcon onClick={onCancel}>
              <IconX />
            </ActionIcon>
          </>
        )}
      </Group>
    </Card>
  );
}

function shootOutToScore({ score, priority }) {
  let fullScore = score;
  if (priority != null) {
    fullScore += priority ? "+" : "-";
  }
  return fullScore;
}

function scoreToShootOut(id, fullScore) {
  let score;
  let priority;
  if (fullScore.endsWith("+")) {
    score = fullScore.slice(0, -1);
    priority = true;
  } else if (fullScore.endsWith("-")) {
    score = fullScore.slice(0, -1);
    priority = false;
  } else {
    score = fullScore;
    priority = null;
  }

  return {
    id,
    score,
    priority,
  };
}
