
type AgentState = 'ACTIVE' | 'INACTIVE' | 'MAINTENANCE' | 'ENROLLING';

interface Agent {
  id: string;
  display_name: string;
  registered_ip: string;
  metadata?: AgentMetadata;
  status: AgentState;
  last_handshake?: number;
}

interface AgentMetadata {
  system_os?: string;
  cpu_model?: string;
  total_memory: number;
  total_disk: number;
  partitions: AgentPartitions[];
}

interface AgentPartitions {
  mount_point: string;
  total_size: number;
  raid: boolean;
  raid_level?: string;
}


type AgentCpuMetric = {
  timestamp: number;
  value: number
}

export { Agent, AgentMetadata, AgentPartitions, AgentState, AgentCpuMetric };
