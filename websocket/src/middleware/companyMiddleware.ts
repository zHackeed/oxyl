import type { Socket } from "socket.io";
import { Middleware } from "../utils/middleware.js";
import { getPermissionsInCompany } from "../storage/company.js";
import { CompanyPermission, hasPermission } from "../types/permissions.js";
import { logger } from "../utils/logConfig.js";
import type { RoomType } from "../types/socket.js";

export class CompanyMiddleware extends Middleware {

  constructor(
    _eventName: string
  ) {
    super(_eventName);
  }

  override async validate(
    userConnection: Socket,
    type: RoomType,
    id: string,
  ): Promise<void | Error> {
    if (type !== "company") return; // not our concern, let it pass

    if (!id) return new Error("Company ID is required");

    const userId = userConnection.data.userId;
    if (!userId) return new Error("User ID is required");

    const permissions = await getPermissionsInCompany(id, userId);
    if (!permissions) return new Error("User does not have permission in this company");

    if (!hasPermission(permissions, CompanyPermission.View))
      return new Error("User does not have permission to view this company");

    logger.info(`User ${userId} has permission to view company ${id}`);
  }

}