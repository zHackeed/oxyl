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

type NotificationType =
  | 'COMPANY_SETTING_UPDATE'
  | 'COMPANY_MEMBER_UPDATE'
  | 'AGENT_STATUS_UPDATE'
  | 'AGENT_CPU_USAGE_THRESHOLD'
  | 'AGENT_MEMORY_USAGE_THRESHOLD'
  | 'AGENT_DISK_USAGE_THRESHOLD'
  | 'AGENT_DISK_HEALTH_THRESHOLD'
  | 'AGENT_NETWORK_USAGE_THRESHOLD';

const notificationTypeLabels: Record<NotificationType, string> = {
  COMPANY_SETTING_UPDATE: 'Actualización de configuración de empresa',
  COMPANY_MEMBER_UPDATE: 'Actualización de miembro de empresa',
  AGENT_STATUS_UPDATE: 'Actualización de estado del agente',
  AGENT_CPU_USAGE_THRESHOLD: 'Umbral de uso de CPU',
  AGENT_MEMORY_USAGE_THRESHOLD: 'Umbral de uso de memoria',
  AGENT_DISK_USAGE_THRESHOLD: 'Umbral de uso de disco',
  AGENT_DISK_HEALTH_THRESHOLD: 'Umbral de salud del disco',
  AGENT_NETWORK_USAGE_THRESHOLD: 'Umbral de uso de red',
};

interface AgentNotificationLog {
  identifier: string;
  agent: string;
  trigger_reason: NotificationType;
  trigger_value: string;
  ack: boolean;
  failed: boolean;
  sent_at: string;
}

type AgentCpuMetric = {
  timestamp: number;
  value: number;
};

export {
  Agent,
  AgentMetadata,
  AgentPartitions,
  AgentState,
  AgentCpuMetric,
  AgentNotificationLog,
  NotificationType,
  notificationTypeLabels,
};
