declare namespace NodeJS {
  interface ProcessEnv {
    REDIS_HOST: string;
    REDIS_PORT: number;
    REDIS_PASS: string;
    REDIS_DB: number;
    WEB_URL: string;
    MAIN_URL: string;
    AUTH_URL: string;
    SESSION_NAME: string;
  }
}
