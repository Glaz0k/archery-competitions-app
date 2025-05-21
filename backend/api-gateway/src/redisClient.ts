import Redis from "ioredis";

const { REDIS_HOST, REDIS_PORT, REDIS_PASS, REDIS_DB } = process.env;

export const redisClient = new Redis({
  host: REDIS_HOST,
  port: REDIS_PORT,
  password: REDIS_PASS,
  db: REDIS_DB,
});

redisClient.on("error", (err) => {
  console.error(`Redis error: ${err}`);
});
