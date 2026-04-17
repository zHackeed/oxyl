
import pqPool from "../db/pqConn.js";
import { logger } from "../utils/logConfig.js";

// -> thonking
export async function getPermissionsInCompany(companyId: string, userId: string): Promise<number | undefined> {
  const client = await pqPool.connect();
  try {
    const result = await client.query(`SELECT permission_bitwise FROM company_members WHERE company_id = $1 AND user_id = $2`, [companyId, userId]);
    
    return result.rows[0]?.permission_bitwise;
  } catch (err) {
    logger.error('Error getting permissions in company', err);
    throw err;
  } finally {
    client.release();
  }
}
