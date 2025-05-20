import { RequestHandler } from "http-proxy-middleware";

const { SESSION_NAME } = process.env;

export const logMiddleware: RequestHandler = async (req, _res, next) => {
  console.log(`Invoked: ${req.method} ${req.url}`);
  console.log(`Agent: ${req.headers["user-agent"]}`);
  console.log(`Session: ${req.cookies[SESSION_NAME] ?? "unauthorized"}`);

  next();
};
