import { Company, CompanyPermission } from '@/lib/api/models/company';
import { createWithEqualityFn } from 'zustand/traditional';

//Todo: rethink if we really need this?

export interface CompanyState {
  selectedCompany: Company | null;
  companyPermission: CompanyPermission[] | null;
  setCompany: (company: Company | null) => void;
  setCompanyPermissions: (permissions: CompanyPermission[] | null) => void;
}

const initalState = {
  selectedCompany: null,
  companyPermission: null,
};

export const useCompanyStore = createWithEqualityFn<CompanyState>((set) => ({
  ...initalState,
  setCompany: (company: Company | null) => {
    set((state) => ({
      ...state,
      selectedCompany: company,
    }));
  },
  setCompanyPermissions: (permissions: CompanyPermission[] | null) => {
    set((state) => ({
      ...state,
      companyPermission: permissions,
    }));
  },
}));
