import apiClient from "../../axios/config";
import { AuthDataSchema } from "./schemas";
import type { AuthData, Credentials } from "./types";

export const authApi = {
  adminSignIn: async (credentials: Credentials): Promise<AuthData> => {
    const response = await apiClient.post("/auth/sign_in", credentials);
    const validatedResponse = AuthDataSchema.parse(response.data);
    if (validatedResponse.role !== "admin") {
      throw new Error(`Invalid authority: ${validatedResponse.role}`);
    }
    return validatedResponse;
  },
  logout: async (): Promise<void> => {
    await apiClient.post("/auth/logout");
  },
};
