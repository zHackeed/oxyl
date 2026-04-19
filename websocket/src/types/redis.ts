export type AgentState = "ACTIVE" | "INACTIVE" | "ENROLLING" | "MAINTENANCE";

export const RedisChannels = {
  CompanyStartedListening: "company:state:listening",
  CompanyStoppedListening: "company:state:stopped-listening",

  CompanyAgentCreated: "company:agent:creation",
  CompanyAgentRemoved: "company:agent:removed",
  CompanyAgentStateUpdated: "company:agent:state:updated",
} as const;


export interface CompanyStartedListeningMessage {
  company: string;
}

export interface CompanyStoppedListeningMessage {
  company: string;
}


// ! -> API Events

export interface AgentCreateMessage {
  company_id: string;
  agent_id: string;
  display_name: string;
  state: AgentState;
  registered_ip: string;
}

export interface AgentRemovedMessage {
  company_id: string;
  agent_id: string;
}

// ! -> Ingress events
export interface AgentStateUpdateMessage {
  company: string;
  agent: string;
  state: AgentState;
};

