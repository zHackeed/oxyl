import pqPool from "../db/pqConn.js";
import type { Company } from "../types/company.js";
import type { CompanyPermission } from "../types/permissions.js";
import { logger } from "../utils/logConfig.js";
import { Cache } from "../utils/cache.js";

const companyPermissionCache = new Cache<string, CompanyPermission>(1000);

export async function getCompany(
  companyId: string,
): Promise<Company | undefined> {
  const client = await pqPool.connect();
  try {
    const result = await client.query(
      `SELECT id, display_name, holder, limit_agents FROM companies WHERE id = $1 AND enabled = true`,
      [companyId],
    );

    if (!result.rows[0]) {
      return undefined;
    }

    return {
      id: result.rows[0].id,
      display_name: result.rows[0].display_name,
      //agent_count: result.rows[0].agent_count,
      agent_count: 0,
      limit_agents: result.rows[0].limit_agents,
      holder: result.rows[0].holder,
    };
  } catch (err) {
    logger.error("Error getting company", err);
    throw err;
  } finally {
    client.release();
  }
}

export async function getPermissionsInCompany(
  companyId: string,
  userId: string,
): Promise<CompanyPermission | undefined> {
  const cachedValue = companyPermissionCache.get(`${companyId}:${userId}`);
  if (cachedValue) {
    return cachedValue;
  }

  const client = await pqPool.connect();
  try {
    const result = await client.query(
      `SELECT permission_bitwise FROM company_members WHERE company_id = $1 AND user_id = $2`,
      [companyId, userId],
    );

    const value = result.rows[0]?.permission_bitwise;
    if (value) {
      companyPermissionCache.set(`${companyId}:${userId}`, value, 60000);
    }

    return value;
  } catch (err) {
    logger.error("Error getting permissions in company", err);
    throw err;
  } finally {
    client.release();
  }
}
