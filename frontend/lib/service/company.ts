
import { Caller } from "../api/api";
import { ActiveCompanyThreshold, CompanyMember, CompanyThresholdNotificationType } from "../api/models/company";
import { CreateCompanyRequest } from "../api/requests/company";
import { CreateCompanyResponse, GetCompanyPermissionResponse } from "../api/responses/company";
import { createCompanySchemaValidator } from "../validators/company";

export const companyService = {

  createCompany: async (data: CreateCompanyRequest) : Promise<CreateCompanyResponse> => {
    try {
      const validated = await createCompanySchemaValidator.validate(data)
      if (!validated) {
        throw new Error('Invalid company data')
      }

      const response = await Caller.post<CreateCompanyResponse>('company/create', validated)
    
      if (response.status !== 201) {
        throw new Error('Failed to create company')
      }

      return response.data
    } catch (error) {
      return Promise.reject(error)
    }
  },
  
  getCompanies: async () => {
    const response = await Caller.get('company');
    
    if (response.status !== 200) {
      throw new Error('Failed to fetch companies');
    }
    
    return response.data;
  },
  
  getCompany: async (id: string) => {
    const response = await Caller.get(`/company/${id}`);
    
    if (response.status !== 200) {
      throw new Error('Failed to fetch company');
    }
    
    return response.data;
  },

  getPermission: async (id: string) : Promise<GetCompanyPermissionResponse> =>  {
    const response = await Caller.get(`/company/${id}/permissions/self`);
    
    if (response.status !== 200) {
      throw new Error('Failed to fetch permissions for company', {
        cause: {
          error_code: response.status,
        }
      });
    }
    
    return response.data;
  },

  getActiveCompanyThresholds: async (id: string) : Promise<ActiveCompanyThreshold[]> => {
    const response = await Caller.get(`/company/${id}/thresholds`);
    
    if (response.status !== 200) {
      throw new Error('Failed to fetch active company thresholds');
    }
    
    return response.data;
  },

  updateThreshold: async (id: string, key: CompanyThresholdNotificationType, value: number) => {
    const response = await Caller.patch(`/company/${id}/thresholds`, {
      notification_type: key,
      threshold: value,
    });
    
    if (response.status !== 200) {
      throw new Error('Failed to update threshold');
    }
    
    return response.data;
  },

  getMembers: async (id: string) : Promise<CompanyMember[]> => {
    const response = await Caller.get(`/company/${id}/member`);
    
    if (response.status !== 200) {
      throw new Error('Failed to fetch members');
    }
    
    return response.data;
  }

}