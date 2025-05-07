import type { PropsWithChildren } from "react";
import { AspectRatio, Card, Center, Title } from "@mantine/core";
import type { Score } from "../../../entities";
import { NO_SCORE_VALUE } from "../../constants";

export interface ScoreCardProps {
  score: Score;
  editing: boolean;
}

export function ScoreCard({ score, editing, children }: PropsWithChildren<ScoreCardProps>) {
  return (
    <AspectRatio ratio={1}>
      <Card p="sm">
        <Center w={40} h={40}>
          {editing ? children : <Title order={2}>{score ?? NO_SCORE_VALUE}</Title>}
        </Center>
      </Card>
    </AspectRatio>
  );
}
