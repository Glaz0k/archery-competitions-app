import { useEffect } from "react";
import { Outlet, useNavigate } from "react-router";
import { Center, Loader } from "@mantine/core";
import { useUser } from "../../../api";

export function ProtectedRoute() {
  const { data: authData, isFetching } = useUser();
  const navigate = useNavigate();

  useEffect(() => {
    if (!isFetching && !authData) {
      navigate("/sign-in");
    }
  }, [authData, isFetching, navigate]);

  if (isFetching || !authData) {
    return (
      <Center h="100vh" flex={1}>
        <Loader size="xl" />
      </Center>
    );
  }

  return <Outlet />;
}
