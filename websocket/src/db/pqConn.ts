import { Pool } from "pg";
import { logger } from "../utils/logConfig.js";

const pqPool: Pool = new Pool({
  host: process.env["TIGERDB_HOST"],
  port: Number(process.env["TIGERDB_PORT"]),
  user: process.env["TIGERDB_USER"],
  password: process.env["TIGERDB_PASS"],
  database: process.env["TIGERDB_DB"],
});

async function verifyConnection() {
  try {
    const client = await pqPool.connect();
    logger.info("Database connection verified");
    client.release();
  } catch (error) {
    logger.error("Database connection verification failed", error);
    process.exit(1);
  }
}

verifyConnection();

export default pqPool;
