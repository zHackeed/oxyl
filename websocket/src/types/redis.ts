import type {
  AgentGeneralMetric,
  AgentMountPointMetric,
  AgentNetworkMetric,
  AgentPhysicalDiskMetric,
} from "./metrics.js";

export type AgentState = "ACTIVE" | "INACTIVE" | "ENROLLING" | "MAINTENANCE";

export const RedisChannels = {
  UserInvalidation: "user:invalidate",

  CompanyMemberAdd: "company:member:add",
  CompanyMemberRemove: "company:member:remove",
  CompanyDeletion: "company:deletion",

  CompanyAgentCreated: "company:agent:creation",
  CompanyAgentRemoved: "company:agent:removed",
  CompanyAgentStateUpdated: "company:agent:state:updated",

  AgentStateUpdate: "agent:state:update",

  AgentStartedListening: "agent:viewer:listening",
  AgentMetricAppend: "agent:viewer:metric:append",
  AgentStoppedListening: "agent:viewer:deafen", // deafened = not listening, it's a joke
} as const;

export interface AgentStartedListeningMessage {
  agent: string;
}

export interface AgentStoppedListeningMessage {
  agent: string;
}

// ! -> API Events

export interface UserInvalidationMessage {
  user_id: string;
}

export interface CompanyMemberInvitationMessage {
  company_id: string;
  user_id: string;
  permissions: number;
}

export interface CompanyMemberRemovalMessage {
  company_id: string;
  user_id: string;
}

export interface CompanyDeletionMessage {
  company_id: string;
  affected: string[];
}

export interface AgentCreateMessage {
  company_id: string;
  agent_id: string;
  display_name: string;
  state: AgentState;
  registered_ip: string;
}

export interface AgentMetricAppendMessage {
  agent_id: string;
  general: AgentGeneralMetric;
  network: AgentNetworkMetric[];
  mount: AgentMountPointMetric[];
  physical_disk: AgentPhysicalDiskMetric[];
}

export interface AgentRemovedMessage {
  company_id: string;
  agent_id: string;
}

// ! -> Ingress events (from API/other services)
export interface AgentStateUpdateMessage {
  company_holder: string;
  agent_id: string;
  status: AgentState;
}
