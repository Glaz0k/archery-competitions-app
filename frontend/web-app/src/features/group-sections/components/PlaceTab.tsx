import { useState } from "react";
import { z } from "zod";
import { TextInput, useMantineTheme } from "@mantine/core";
import { useForm, zodResolver } from "@mantine/form";
import {
  useCompletePlaceRange,
  useEditPlaceRange,
  useEditShootOut,
  type ShootOutEdit,
} from "../../../api";
import { RangeType, type Range, type ShootOut, type SparringPlace } from "../../../entities";
import { RangeCard, RangeSection, ScoreCard, ShotsFormControls } from "../../range-section";

export interface PlaceTabProps {
  place: SparringPlace;
}

export function PlaceTab({ place: { id, rangeGroup, shootOut } }: PlaceTabProps) {
  return (
    <>
      {[...rangeGroup.ranges]
        .sort(({ ordinal: a }, { ordinal: b }) => a - b)
        .map((range) => (
          <PlaceRangeCard
            key={`${id}$range${range.ordinal}`}
            placeId={id}
            range={range}
            rangeSize={rangeGroup.rangeSize}
            rangeType={rangeGroup.type}
          />
        ))}
      {shootOut && <ShootOutCard placeId={id} shootOut={shootOut} type={rangeGroup.type} />}
    </>
  );
}

interface PlaceRangeCardProps {
  placeId: number;
  range: Range;
  rangeSize: number;
  rangeType: RangeType;
}

function PlaceRangeCard({ placeId, range, rangeSize, rangeType }: PlaceRangeCardProps) {
  const { mutate: editRange, isPending: isRangeEditing } = useEditPlaceRange();
  const { mutate: completeRange, isPending: isRangeCompleting } = useCompletePlaceRange();

  const shotOrdinals = [...Array(rangeSize).keys()].map((i) => i + 1);
  const shots = shotOrdinals.map((ord) => {
    let shot = range.shots?.find((shot) => shot.ordinal === ord);
    if (!shot) {
      shot = {
        ordinal: ord,
        score: null,
      };
    }
    return shot;
  });

  return (
    <RangeCard
      active={range.isActive}
      title={`Серия ${range.ordinal}`}
      loading={isRangeEditing || isRangeCompleting}
    >
      <RangeSection
        shots={shots}
        type={rangeType}
        editFn={(editedShots) =>
          editRange([placeId, { ordinal: range.ordinal, shots: editedShots }])
        }
        completeFn={() => completeRange([placeId, range.ordinal])}
        active={range.isActive}
      />
    </RangeCard>
  );
}

interface ShootOutCardProps {
  placeId: number;
  shootOut: ShootOut;
  type: RangeType;
}

function ShootOutCard({ placeId, shootOut, type }: ShootOutCardProps) {
  const theme = useMantineTheme();
  const shootOutRg =
    type === RangeType.ONE_TEN ? /^(M|[1-9]|10|X)([+-])?$/ : /^(M|[6-9]|10|X)([+-])?$/;
  const shootOutForm = useShootOutForm(shootOut, shootOutRg);
  const [isShootOutEditing, setShootOutEditing] = useState<boolean>(false);

  const { mutate: submitShootOut, isPending: isShootOutSubmitting } = useEditShootOut(() => {
    setShootOutEditing(false);
  });

  return (
    <RangeCard active={shootOut.score === null} title="Перестрелка" loading={isShootOutSubmitting}>
      <form
        onSubmit={shootOutForm.onSubmit((edited) => {
          submitShootOut([placeId, edited]);
        })}
      >
        <ScoreCard score={shootOutToFullScore(shootOut)} editing={isShootOutEditing}>
          <TextInput
            p={0}
            variant="unstyled"
            styles={{
              input: {
                font: theme.headings.fontFamily,
                fontSize: theme.headings.sizes.h2.fontSize,
                fontWeight: theme.headings.sizes.h2.fontWeight,
                color:
                  shootOutForm.errors?.fullScore === undefined ? theme.white : theme.colors.red[6],
              },
            }}
            key={shootOutForm.key("fullScore")}
            {...shootOutForm.getInputProps("fullScore", { withError: false })}
          />
        </ScoreCard>
        <ShotsFormControls
          editing={isShootOutEditing}
          onEdit={() => {
            shootOutForm.setValues({
              fullScore: shootOutToFullScore(shootOut),
            });
            setShootOutEditing(true);
          }}
          onCancelEdit={() => {
            setShootOutEditing(false);
            shootOutForm.reset();
          }}
        />
      </form>
    </RangeCard>
  );
}

const useShootOutForm = (initial: ShootOut, shootOutRegex: RegExp) => {
  const ShootOutFormSchema = z.object({
    fullScore: z.union([
      z.string().trim().max(0).nullable(),
      z.string().trim().regex(shootOutRegex),
    ]),
  });
  type ShootOutFormValues = z.infer<typeof ShootOutFormSchema>;

  return useForm<ShootOutFormValues, (values: ShootOutFormValues) => ShootOutEdit>({
    mode: "uncontrolled",
    initialValues: {
      fullScore: shootOutToFullScore(initial),
    },
    validateInputOnChange: true,
    validate: zodResolver(ShootOutFormSchema),
    transformValues: ({ fullScore }) => {
      if (!fullScore) {
        return {
          score: null,
          priority: null,
        };
      }
      let score: string | null = null;
      let priority: boolean | null = null;
      if (fullScore.endsWith("+")) {
        score = fullScore.slice(0, -1);
        priority = true;
      } else if (fullScore.endsWith("-")) {
        score = fullScore.slice(0, -1);
        priority = false;
      } else if (fullScore !== "") {
        score = fullScore;
      }
      return {
        score,
        priority,
      };
    },
  });
};

const shootOutToFullScore = (shootOut: ShootOut): string => {
  const value = shootOut.score;
  if (value === null) {
    return "";
  }
  if (shootOut.priority == null) {
    return value;
  }
  return shootOut.priority ? `${value}+` : `${value}-`;
};
