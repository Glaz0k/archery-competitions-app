import express from "express";
import cookieParser from "cookie-parser";
import cors, { CorsOptions } from "cors";
import { sessionMiddleware } from "./proxy/sessionMiddleware";
import { config } from "dotenv";
import { handleLogout } from "./endpoints/handleLogout";
import { jwtMiddleware } from "./jwtMiddleware";
import { authMiddleware } from "./proxy/authMiddleware";
import { mainMiddleware } from "./proxy/mainMiddleware";
import { logMiddleware } from "./logMiddleware";
config();

const { WEB_URL, MAIN_URL, AUTH_URL } = process.env;

const corsOptions: CorsOptions = {
  origin: WEB_URL,
  credentials: true,
};

const app = express();

app.use(cors(corsOptions));
app.use(cookieParser());

app.use("/", logMiddleware);
app.use("/api/auth/sign_up", sessionMiddleware);
app.use("/api/auth/sign_in", sessionMiddleware);
app.post("/api/auth/logout", handleLogout);
app.use("/api/auth", jwtMiddleware, authMiddleware);
app.use("/api/", jwtMiddleware, mainMiddleware);

app
  .listen(3000, async () => {
    console.log(`Gateway running on port 3000`);
    console.log(`Main service URL: ${MAIN_URL}`);
    console.log(`Auth service URL: ${AUTH_URL}`);
    console.log(`Web URL: ${WEB_URL}`);
  })
  .on("error", (err) => {
    console.error(`Failed to start gateway: ${err}`);
    process.exit(1);
  });
