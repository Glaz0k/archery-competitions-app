import { createProxyMiddleware } from "http-proxy-middleware";

const { MAIN_URL } = process.env;

export const mainMiddleware = createProxyMiddleware({
  target: MAIN_URL,
  changeOrigin: true,
});
