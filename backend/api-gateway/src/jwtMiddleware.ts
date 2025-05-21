import { RequestHandler } from "http-proxy-middleware";
import { redisClient } from "./redisClient";

const { SESSION_NAME } = process.env;

export const jwtMiddleware: RequestHandler = async (req, res, next) => {
  const cookie = req.cookies[SESSION_NAME];

  if (cookie) {
    try {
      const token = await redisClient.get(cookie);
      if (token) {
        req.headers["authorization"] = `Bearer ${token}`;
      }
    } catch (err) {
      console.error("Jwt middleware error:", err);
      res.status(500).json({ error: "GATEWAY JWT ERROR" });
      return;
    }
  }

  next();
};
