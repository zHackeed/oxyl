import { Agent } from "../api/models/agent";

export type RoomType = "company" | "agent";

export interface UserSocketReq {
  join: (type: RoomType, id: string) => void;
  leave: (type: RoomType, id: string) => void;
}

export interface CompanyUpdateActions {
  "agent:creation": (agent: Agent) => void;
  //"agent:update": (agentId: string, state: AgentStatus) => void;
  "agent:deletion": (agent: string) => void;

  "user:add": (userId: string) => void;
  "user:remove": (userId: string) => void;
}

// :thinking:  
export interface SocketMetadata {
  userId: string;
  jti: string;
}

// !! --- company state
export const AgentUpdateAction = {
  CREATED: "created",
  UPDATED: "updated",
  DELETED: "deleted",
} as const;