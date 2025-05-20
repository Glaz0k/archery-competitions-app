import type { PropsWithChildren } from "react";
import { Center, Loader } from "@mantine/core";
import { CenterCard } from "../cards/CenterCard";

export interface PageLoaderProps {
  loading: boolean;
  error: boolean;
}

export function PageLoader({ loading, error, children }: PropsWithChildren<PageLoaderProps>) {
  if (loading) {
    return (
      <Center flex={1}>
        <Loader size="xl" />
      </Center>
    );
  }
  if (error) {
    return <CenterCard label="Произошла ошибка" />;
  }
  return children;
}
