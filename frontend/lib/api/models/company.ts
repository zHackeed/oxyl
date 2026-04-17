import { GetThemeValueForKey } from 'tamagui';
import { UserResumed } from './user';

export interface Company {
  id: string;
  display_name: string;
  holder: string; // Company  holder <-> user
  limit_nodes: number;
}

export interface ActiveCompanyThreshold {
  threshold_id: CompanyThresholdNotificationType;
  value: number;
}

export const CompanyThresholdNotificationType = {
  //AgentStatusUpdate: "AGENT_STATUS_UPDATE",
  AgentCpuUsageThreshold: 'AGENT_CPU_USAGE_THRESHOLD',
  AgentMemoryUsageThreshold: 'AGENT_MEMORY_USAGE_THRESHOLD',
  AgentDiskUsageThreshold: 'AGENT_DISK_USAGE_THRESHOLD',
  AgentDiskHealthThreshold: 'AGENT_DISK_HEALTH_THRESHOLD',
  AgentNetworkThreshold: 'AGENT_NETWORK_THRESHOLD',
} as const;

export type CompanyThresholdNotificationType =
  (typeof CompanyThresholdNotificationType)[keyof typeof CompanyThresholdNotificationType];

export const CompanyPermission = {
  View: 'view',
  ManageAgents: 'manage_agents',
  ManageCompany: 'manage_company',
  ManageMembers: 'manage_members',
  ManageThresholds: 'manage_thresholds',
  ManageWebhooks: 'manage_webhooks',
  Owner: 'owner',
} as const;

export type CompanyPermission = (typeof CompanyPermission)[keyof typeof CompanyPermission];

export function hasPermission(
  permissions: CompanyPermission[],
  required: CompanyPermission
): boolean {
  return permissions.includes(required);
}

// ------------------------

export interface CompanyMember {
  user: UserResumed;
  permissions: CompanyPermission[];
  created_at: Date; //
}

/// ------------------------

export interface ThresholdMetadata {
  label: string;
  description: string;
  color: GetThemeValueForKey<`color`>;
}

export const ThresholdMetadataMap: Record<CompanyThresholdNotificationType, ThresholdMetadata> = {
  [CompanyThresholdNotificationType.AgentCpuUsageThreshold]: {
    label: 'Uso de CPU',
    description: 'Limite de uso aceptable medio de CPU',
    color: '$red9',
  },
  [CompanyThresholdNotificationType.AgentMemoryUsageThreshold]: {
    label: 'Uso de memoria',
    description: 'Limite de uso aceptable medio de memoria',
    color: '$yellow9',
  },
  [CompanyThresholdNotificationType.AgentDiskUsageThreshold]: {
    label: 'Uso de disco',
    description: 'Limite de uso aceptable medio de disco',
    color: '$blue9',
  },
  [CompanyThresholdNotificationType.AgentDiskHealthThreshold]: {
    label: 'Salud del disco',
    description: 'Limite de salud aceptable del disco',
    color: '$orange9',
  },
  [CompanyThresholdNotificationType.AgentNetworkThreshold]: {
    label: 'Uso de red',
    description: 'Limite de uso aceptable medio de red (TX/RX)',
    color: '$blue9',
  },
} as const;

// The ugliest stuff you are going to see in your life
export function isAdmin(member: CompanyMember): boolean {
  return (
    member.permissions.includes(CompanyPermission.Owner) ||
    (member.permissions.includes(CompanyPermission.ManageCompany) &&
      member.permissions.includes(CompanyPermission.ManageMembers) &&
      member.permissions.includes(CompanyPermission.ManageThresholds) &&
      member.permissions.includes(CompanyPermission.ManageWebhooks))
  );
}
