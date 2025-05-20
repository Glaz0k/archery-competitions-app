import { z } from "zod";

export const CredentialsSchema = z.object({
  login: z.string().regex(/^[a-zA-Z0-9._-]{6,20}$/, { message: "Неверное имя пользователя" }),
  password: z.string().regex(/^[a-zA-Z0-9._-]{6,20}$/, { message: "Неверный пароль" }),
});

export const AuthDataSchema = z.object({
  user_id: z.number(),
  role: z.union([z.literal("admin"), z.literal("user")]),
});
