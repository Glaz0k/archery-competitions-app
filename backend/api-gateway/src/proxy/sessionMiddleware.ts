import { createProxyMiddleware, RequestHandler } from "http-proxy-middleware";
import { redisClient } from "../redisClient";

const { AUTH_URL, SESSION_NAME } = process.env;

interface SignInProxyResponse {
  auth_data: {
    user_id: number;
    role: string;
  };
  token: string;
}

export const sessionMiddleware: RequestHandler = createProxyMiddleware({
  target: AUTH_URL,
  changeOrigin: true,
  selfHandleResponse: true,
  onProxyRes: async (proxyRes, req, res) => {
    let body: string;
    try {
      body = await new Promise((resolve) => {
        let data = "";
        proxyRes.on("data", (chunk) => {
          data += chunk;
        });
        proxyRes.on("end", () => {
          resolve(data);
        });
      });
    } catch (err) {
      console.error("Response body error:", err);
      res.status(500).json({ error: "GATEWAY SESSION MIDDLEWARE ERROR" });
      return;
    }

    if (proxyRes.statusCode !== 200) {
      res.status(proxyRes.statusCode ?? 500).json(JSON.parse(body));
      return;
    }

    try {
      const { auth_data, token }: SignInProxyResponse = JSON.parse(body);

      const cookie = `${auth_data.user_id}-${auth_data.role}`;

      await redisClient.set(cookie, token);

      res
        .cookie(SESSION_NAME, cookie, {
          httpOnly: true,
          secure: false,
          sameSite: "lax",
        })
        .json(auth_data);
    } catch (err) {
      console.error("Session middleware error:", err);
      res.status(500).json({ error: "GATEWAY SESSION MIDDLEWARE ERROR" });
    }
  },
});
