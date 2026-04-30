import type { Server, Socket } from "socket.io";
import type { RedisMessenger } from "../db/redisConn.js";
import { RedisChannels, type AgentStateUpdateMessage, type AgentCreateMessage, type AgentRemovedMessage, type CompanyDeletionMessage, type CompanyMemberInvitationMessage, type CompanyMemberRemovalMessage } from "../types/redis.js";
import { getCompany } from "../storage/company.js";
import type { CompanyUpdateActions, SocketMetadata, UserSocketReq } from "../types/socket.js";
import { logger } from "../utils/logConfig.js";
import type { Agent } from "../types/company.js";

export class CompanyService {

  constructor(
    private readonly _messenger: RedisMessenger,
    private readonly _io: Server<
      UserSocketReq,
      CompanyUpdateActions,
      SocketMetadata
    >
  ) {

    this._messenger.subscribe<CompanyMemberInvitationMessage>(
      RedisChannels.CompanyMemberAdd,
      this.handleCompanyMemberAdd.bind(this)
    );

    this._messenger.subscribe<CompanyMemberRemovalMessage>(
      RedisChannels.CompanyMemberRemove,
      this.handleCompanyMemberRemove.bind(this)
    );

    this._messenger.subscribe<CompanyDeletionMessage>(
      RedisChannels.CompanyDeletion,
      this.handleCompanyDeletion.bind(this)
    )

    this._messenger.subscribe<AgentCreateMessage>(
      RedisChannels.CompanyAgentCreated,
      this.handleAgentCreation.bind(this)
    );
    this._messenger.subscribe<AgentStateUpdateMessage>(
      RedisChannels.AgentStateUpdate,
      this.handleAgentStateUpdate.bind(this)
    );
    this._messenger.subscribe<AgentRemovedMessage>(
      RedisChannels.CompanyAgentRemoved,
      this.handleAgentDeletion.bind(this)
    );
  }

  async addListener(companyId: string, user: Socket) {
    const companyRoom = this._ch(companyId)

    // const connectedSockets = await this._io.in(companyRoom).fetchSockets()
    // if (connectedSockets.length === 0) {
    //   await this._messenger.publish(RedisChannels.CompanyStartedListening, {
    //     company: companyId,
    //   } as CompanyStartedListeningMessage);
    // }

    if (!user.rooms.has(companyRoom)) {
      user.join(companyRoom)
      logger.info("user has joined the company room", companyId)
    }
  }

  async removeListener(companyId: string, user: Socket) {
    const companyRoom = this._ch(companyId)
    // const connectedSockets = await this._io.in(companyRoom).fetchSockets()

    // // If the user is not in the room, we can't remove them. Why would this happen normally? Maybe the user disconnected before the removeListener was called?
    // if (connectedSockets.find((socket) => socket.id === user.id)) return;

    // if ((connectedSockets.length - 1) < 1) {
    //   await this._messenger.publish(RedisChannels.CompanyStoppedListening, {
    //     company: companyId,
    //   } as CompanyStartedListeningMessage);
    // }

    if (user.rooms.has(companyRoom)) {
      user.leave(companyRoom)
      logger.info("user has left the company room", companyId)
    }
  }

  // * redis incoming messages  services <-> websocket
  // todo: think about this logic
  private async handleCompanyMemberAdd(message: CompanyMemberInvitationMessage) {
    const company = await getCompany(message.company_id);
    if (!company) return;

    this._io.in(`user:${message.user_id}`).emit("company:added", company)
    //this._io.in(this._ch(message.company_id)).emit("company:member:added", message.user_id, message.permissions)
  }

  private async handleCompanyMemberRemove(message: CompanyMemberRemovalMessage) {
    this._io.in(`user:${message.user_id}`).emit("company:removed", message.company_id)
    this._io.in(this._ch(message.company_id)).emit("company:member:removed", message.user_id)
  }

  private async handleCompanyDeletion(message: CompanyDeletionMessage) {
    this._io.in(this._ch(message.company_id)).emit('company:removed', message.company_id)

    for (const userId of message.affected) {
      this._io.in(`user:${userId}`).emit("company:removed", message.company_id)
    }
  }

  private async handleAgentCreation(message: AgentCreateMessage) {
    const agent: Agent = {
      id: message.agent_id,
      display_name: message.display_name,
      registered_ip: message.registered_ip,
      status: message.state,
    }
    
    this._io.in(this._ch(message.company_id)).emit("company:agent:creation", agent)
  }

  private async handleAgentStateUpdate(message: AgentStateUpdateMessage) {
    this._io.in(this._ch(message.company_holder)).emit(`company:agent:update`, message.agent_id, message.status)
  }

  private async handleAgentDeletion(message: AgentRemovedMessage) {
    this._io.in(this._ch(message.company_id)).emit(`company:agent:deletion`, message.agent_id)
  }

  private _ch(companyId: string) {
    return `company:${companyId}`;
  }
}