export interface AgentMountPointMetric {
  when: string;
  mount_point: string;
  disk_usage: number;
}

export interface AgentPhysicalDiskMetric {
  when: string;
  disk_path: string;
  health_used?: number;
  media_errors_1?: number;
  media_errors_2?: number;
  error_rate?: number;
  pending_sectors?: number;
}

export interface AgentGeneralMetric {
  when: string;
  cpu_usage: number;
  memory_usage: number;
  uptime: number;
}

export interface AgentNetworkMetric {
  when: string;
  interface_name: string;
  rx_bytes: number;
  tx_bytes: number;
  rx_packets: number;
  tx_packets: number;
  rx_rate: number;
  tx_rate: number;
  rx_packet_rate: number;
  tx_packet_rate: number;
}

export interface AgentMetricEntry {
  general: AgentGeneralMetric;
  mount: AgentMountPointMetric[];
  physical_disk: AgentPhysicalDiskMetric[];
  network: AgentNetworkMetric[];
}
