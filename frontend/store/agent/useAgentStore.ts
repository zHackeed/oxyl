import { Agent, AgentState } from "@/lib/api/models/agent";
import { createWithEqualityFn } from "zustand/traditional";


export interface AgentStoreProps {
  agent: Agent | null;
  status: AgentState;
  setAgent: (agent: Agent | null) => void;
  setStatus: (status: AgentState) => void;
}

export const useAgentStore = createWithEqualityFn<AgentStoreProps>()((set) => ({
  agent: null,
  status: 'INACTIVE',
  setAgent: (agent: Agent | null) => {
    set({ agent, status: agent?.status || 'INACTIVE' });
  },
  setStatus: (status: AgentState) => {
    set({ status });
  },
}));