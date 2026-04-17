interface Agent {
  id: string;
  display_name: string;
  registered_ip: string;
  metadata?: AgentMetadata;
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

export { Agent, AgentMetadata, AgentPartitions };
