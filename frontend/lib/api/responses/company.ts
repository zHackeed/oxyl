import { Company, CompanyPermission } from "../models/company";

interface CreateCompanyResponse {
  company: Company;
  permissions: CompanyPermission[];
}

interface GetCompanyPermissionResponse {
  user_id: string;
  permissions: CompanyPermission[];
}

export {
  CreateCompanyResponse,
  GetCompanyPermissionResponse
}