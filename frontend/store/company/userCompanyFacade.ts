import { useCompanyStore } from './useCompanyStore';
import { shallow } from 'zustand/shallow';

export const useCompanyFacade = () => {
  const { activeCompany, permissions, setCompany, setCompanyPermissions } = useCompanyStore(
    (state) => ({
      activeCompany: state.selectedCompany,
      permissions: state.companyPermission,
      setCompany: state.setCompany,
      setCompanyPermissions: state.setCompanyPermissions,
    }),
    shallow
  );

  return {
    activeCompany,
    permissions,
    setCompany,
    setCompanyPermissions,
  };
};
