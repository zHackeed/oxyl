
export interface Company {
  id: string;
  display_name: string;
  holder: string;
  agent_count: number; 
  limit_agents: number;
}

export type AgentState = 'ACTIVE' | 'INACTIVE' | 'MAINTENANCE' | 'ENROLLING';

export interface Agent {
  id: string;
  display_name: string;
  registered_ip: string;
  status: AgentState;
}