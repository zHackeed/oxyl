import type { Server, Socket } from "socket.io";
import type { RedisMessenger } from "../db/redisConn.js";
import { RedisChannels, type AgentStateUpdateMessage, type CompanyStartedListeningMessage, type AgentCreateMessage, type AgentRemovedMessage } from "../types/redis.js";

export class CompanyService {

  constructor(
    private readonly _messenger: RedisMessenger,
    private readonly _io: Server
  ) {
    this._messenger.subscribe<AgentCreateMessage>(
      RedisChannels.CompanyAgentCreated,
      this.handleAgentCreation.bind(this)
    );
    this._messenger.subscribe<AgentStateUpdateMessage>(
      RedisChannels.CompanyAgentStateUpdated,
      this.handleAgentStateUpdate.bind(this)
    );
    this._messenger.subscribe<AgentRemovedMessage>(
      RedisChannels.CompanyAgentRemoved,
      this.handleAgentDeletion.bind(this)
    );
  }

  async addListener(companyId: string, user: Socket) {
    const companyRoom = this.ch(companyId)

    const connectedSockets = await this._io.in(companyRoom).fetchSockets()
    if (connectedSockets.length === 0) {
      await this._messenger.publish(RedisChannels.CompanyStartedListening, {
        company: companyId,
      } as CompanyStartedListeningMessage);
    }

    user.join(`company:${companyId}`)
    console.info("user has joined the company room", companyId, "count", connectedSockets.length + 1)
  }

  async removeListener(companyId: string, user: Socket) {
    const companyRoom = this.ch(companyId)
    const connectedSockets = await this._io.in(companyRoom).fetchSockets()

    if ((connectedSockets.length - 1) < 1) {
      await this._messenger.publish(RedisChannels.CompanyStoppedListening, {
        company: companyId,
      } as CompanyStartedListeningMessage);
    } 

    user.leave(`company:${companyId}`)
    console.info("user has left the company room", companyId, "count", connectedSockets.length - 1, "notified", (connectedSockets.length - 1) === 0)
  }

  // * redis incoming messages
  private async handleAgentCreation(message: AgentCreateMessage) {
    console.log("agent created", message)
    this._io.in(this.ch(message.company_id)).emit("agent:creation", message)
  }

  private async handleAgentStateUpdate(message: AgentStateUpdateMessage) {
    console.log("agent state updated", message)
    this._io.in(this.ch(message.company)).emit(`agent:update`, message.agent, message.state)
  }

  private async handleAgentDeletion(message: AgentRemovedMessage) {
    console.log("agent deleted", message)
    this._io.in(this.ch(message.company_id)).emit(`agent:deletion`, message.agent_id)
  }

  private ch(companyId: string) {
    return `company:${companyId}`;
  }

}