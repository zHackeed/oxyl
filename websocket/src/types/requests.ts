export type JoinLocation = 'agent' | 'company'
 
export interface JoinRequest {
  location: JoinLocation;
  companyId?: string;
  agentId?: string;
}
