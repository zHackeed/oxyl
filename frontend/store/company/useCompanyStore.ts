import { Company, CompanyPermission } from '@/lib/api/models/company';
import { createWithEqualityFn } from 'zustand/traditional';

//Todo: rethink if we really need this?

export interface CompanyState {
  company: Company | null;
  permissions: CompanyPermission[] | null;
  setCompany: (company: Company | null) => void;
  setPermissions: (permissions: CompanyPermission[] | null) => void;
}

const initialState = {
  company: null,
  permissions: null,
};

export const useCompanyStore = createWithEqualityFn<CompanyState>((set) => ({
  ...initialState,
  setCompany: (company: Company | null) => {
    set((state) => ({
      ...state,
      company: company,
    }));
  },
  setPermissions: (permissions: CompanyPermission[] | null) => {
    set((state) => ({
      ...state,
      permissions: permissions,
    }));
  },
}));
