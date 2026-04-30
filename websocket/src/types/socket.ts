import type { Company } from "./company.js";
import type { CompanyPermission } from "./permissions.js";
import type { AgentState, Agent } from "./company.js";
import type { AgentMetricEntry } from "./metrics.js";

export type RoomType = "company" | "agent";

export interface UserSocketReq {
  join: (type: RoomType, id: string) => void;
  leave: (type: RoomType, id: string) => void;
}

export interface CompanyUpdateActions {
  // ----------> Company wide related events
  "company:agent:creation": (agent: Agent) => void;
  "company:agent:update": (agentId: string, state: AgentState) => void;
  "company:agent:deletion": (agentId: string) => void;

  // ----------> Company Member related events
  "company:member:added": (userId: string, permissions: CompanyPermission[]) => void;
  "company:member:removed": (userId: string) => void;

  // ----------> Personal company events
  "company:added": (company: Company) => void;
  "company:removed": (companyId: string) => void;

  // ----------> Agent room related events
  "agent:state:update": (state: AgentState) => void;
  "agent:metric:append": (metric: AgentMetricEntry) => void;
}

// * :thinking:
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