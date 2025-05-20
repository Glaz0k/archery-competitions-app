import { createProxyMiddleware } from "http-proxy-middleware";

const { AUTH_URL } = process.env;

export const authMiddleware = createProxyMiddleware({
  target: AUTH_URL,
  changeOrigin: true,
});
