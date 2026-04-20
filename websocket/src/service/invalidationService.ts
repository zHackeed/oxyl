import type { Server } from "socket.io";
import type { RedisMessenger } from "../db/redisConn.js";
import { RedisChannels, type UserInvalidationMessage } from "../types/redis.js";

export class UserInvalidationService {

  constructor(
    private readonly _io: Server,
    private readonly _messenger: RedisMessenger
  ) {
    this._messenger.subscribe<UserInvalidationMessage>(RedisChannels.UserInvalidation, this.handleUserInvalidation.bind(this))
  }

  private async handleUserInvalidation(message: UserInvalidationMessage) {
    this._io.in(`user:${message.user_id}`).disconnectSockets()
  }
}