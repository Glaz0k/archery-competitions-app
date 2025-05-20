import { createProxyMiddleware } from "http-proxy-middleware";

const { MAIN_URL } = process.env;

export const mainMiddleware = createProxyMiddleware({
  target: MAIN_URL,
  changeOrigin: true,
  onProxyRes: (proxyRes) => {
    delete proxyRes.headers["access-control-allow-origin"];
    delete proxyRes.headers["access-control-allow-credentials"];
  },
});
