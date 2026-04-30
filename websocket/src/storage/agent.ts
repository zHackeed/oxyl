import { logger } from "../utils/logConfig.js";
import pqPool from "../db/pqConn.js";
import { Cache } from "../utils/cache.js";

const agentCompanyMappingCache = new Cache<string, string>(1000);

export async function getAgentCompanyMapping(
  agentId: string,
): Promise<string | undefined> {
  const cachedValue = agentCompanyMappingCache.get(agentId);
  if (cachedValue) {
    return Promise.resolve(cachedValue);
  }

  const client = await pqPool.connect();
  try {
    const result = await client.query(
      `SELECT holder FROM agents WHERE id = $1`,
      [agentId],
    );

    if (result.rows[0]) {
      agentCompanyMappingCache.set(agentId, result.rows[0].holder, 100000);
    }

    return result.rows[0]?.holder;
  } catch (err) {
    logger.error("Error getting agent holder", err);
    throw err;
  } finally {
    client.release();
  }
}
