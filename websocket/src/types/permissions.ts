
export enum CompanyPermission {
  View = 1,              // 00000001
  ManageAgents = 2,      // 00000010
  ManageCompany = 4,     // 00000100
  ManageMembers = 8,     // 00001000
  ManageThresholds = 16, // 00010000
  ManageWebhooks = 32,   // 00100000

  Member = View | ManageAgents,                                                              // 00000011
  Admin = Member | ManageCompany | ManageMembers | ManageThresholds | ManageWebhooks,        // 00111111
  Owner = 255,
}

export function hasPermission(current: CompanyPermission, required: CompanyPermission): boolean {
  return (current & required) !== 0;
}
