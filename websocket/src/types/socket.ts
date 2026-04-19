import type { AgentState } from "./redis.js";

export type RoomType = "company" | "agent";

export interface UserSocketReq {
  join: (type: RoomType, id: string) => void;
  leave: (type: RoomType, id: string) => void;
}

export interface CompanyUpdateActions {
  "agent:creation": (agentId: string) => void;
  "agent:update": (agentId: string, state: AgentState) => void;
  "agent:deletion": (agentId: string) => void;

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