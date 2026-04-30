import { RedisMessenger } from "../db/redisConn.js";
import { Server, Socket } from "socket.io";
import { RedisChannels, type AgentStartedListeningMessage, type AgentStateUpdateMessage, type AgentStoppedListeningMessage, type AgentMetricAppendMessage } from "../types/redis.js";
import { logger } from "../utils/logConfig.js";
import type { AgentMetricEntry } from "../types/metrics.js";

export class AgentService {
  constructor(
    private readonly _messenger: RedisMessenger,
    private readonly _io: Server
  ) {
    this._messenger.subscribe(RedisChannels.AgentStateUpdate, this.handleAgentStateUpdate.bind(this));
    this._messenger.subscribe(RedisChannels.AgentMetricAppend, this.handleAgentMetricAppend.bind(this));
  }

  async addListener(agentId: string, user: Socket) {
    const room = this._room(agentId);
    const connectedSocket = await this._io.in(room).fetchSockets();
    if (connectedSocket.length === 0) {
      await this._messenger.publish<AgentStartedListeningMessage>(RedisChannels.AgentStartedListening, {
        agent: agentId,
      });
    }

    user.join(room);
    logger.info(`User joined agent room: ${room}, notified agent? ${connectedSocket.length === 0 ? "yes" : "no"}`);
  }

  async removeListener(agentId: string, user: Socket) {
    const room = this._room(agentId);

    if (!user.rooms.has(room)) return
    user.leave(room)
    const connectedSocket = await this._io.in(room).fetchSockets();
    if (connectedSocket.length <= 0) {
      await this._messenger.publish<AgentStoppedListeningMessage>(RedisChannels.AgentStoppedListening, {
        agent: agentId,
      });
    }

    logger.info(`User left agent room: ${room}, notified agent? ${connectedSocket.length - 1 <= 0 ? "yes" : "no"}`);
  }

  private handleAgentStateUpdate(message: AgentStateUpdateMessage) {
    console.log("Agent state update:", message);
    this._io.in(this._room(message.agent_id)).emit("agent:state:update", message.status);
  }

  private handleAgentMetricAppend(message: AgentMetricAppendMessage) {
    this._io.in(this._room(message.agent_id)).emit("agent:metric:append", message as AgentMetricEntry);
  }

  private _room(agentId: string) {
    return `agent:${agentId}`;
  }
}