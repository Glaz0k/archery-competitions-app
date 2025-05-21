import type { z } from "zod";
import type { AuthDataSchema, CredentialsSchema } from "./schemas";

export type Credentials = z.infer<typeof CredentialsSchema>;

export type AuthData = z.infer<typeof AuthDataSchema>;
