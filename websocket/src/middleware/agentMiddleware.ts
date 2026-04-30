import type { Socket } from "socket.io";
import { Middleware } from "../utils/middleware.js";
import type { RoomType } from "../types/socket.js";
import { getAgentCompanyMapping } from "../storage/agent.js";
import { getPermissionsInCompany } from "../storage/company.js";
import { CompanyPermission, hasPermission } from "../types/permissions.js";
import { logger } from "../utils/logConfig.js";

export class AgentMiddleware extends Middleware {
  constructor(_eventName: string) {
    super(_eventName);
  }

  override async validate(
    userConnection: Socket,
    type: RoomType,
    id: string,
  ): Promise<void | Error> {
    if (type !== "agent") return; // not our concern, let it pass

    if (!id) return new Error("Agent ID is required");

    const holder = await getAgentCompanyMapping(id);
    if (!holder) return new Error("Agent not found");

    const permissions = await getPermissionsInCompany(
      holder,
      userConnection.data.userId,
    );
    
    if (!permissions)
      return new Error("User does not have permission in this company");

    if (!hasPermission(permissions, CompanyPermission.View)) {
      return new Error("User does not have permission to view this agent");
    }

    logger.info(`User ${userConnection.data.userId} accessed agent ${id}`);
  }
}
