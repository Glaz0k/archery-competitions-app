import { RequestHandler } from "express";
import { redisClient } from "../redisClient";

const { SESSION_NAME } = process.env;

export const handleLogout: RequestHandler = async (req, res) => {
  const cookie = req.cookies[SESSION_NAME!];

  if (cookie) {
    try {
      await redisClient.del(cookie);
    } catch (err) {
      console.error(`Redis error while logout: ${err}`);
      res.status(500).json({ error: "GATEWAY LOGOUT ERROR" });
      return;
    }
  }

  res.clearCookie(SESSION_NAME).status(204).send();
};
