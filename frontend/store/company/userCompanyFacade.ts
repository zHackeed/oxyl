import { useCompanyStore } from './useCompanyStore';
import { shallow } from 'zustand/shallow';

export const useCompanyFacade = () => {
  const { activeCompany, permissions, setCompany, setPermissions } = useCompanyStore(
    (state) => ({
      activeCompany: state.company,
      permissions: state.permissions,
      setCompany: state.setCompany,
      setPermissions: state.setPermissions,
    }),
    shallow
  );

  return {
    activeCompany,
    permissions,
    setCompany,
    setPermissions,
  };
};
