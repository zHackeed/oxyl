import { useAuthStore } from './useAuthStore';
import { shallow } from "zustand/shallow"

// We need this to expose the methods from zustard
export const useAuthFacade = () => {
  const {
    token,
    refreshToken,
    status,
    signIn,
    signOut,
    hydrate,
  } = useAuthStore(
    (state) => ({
      token: state.token,
      refreshToken: state.refreshToken,
      status: state.status,
      signIn: state.signIn,
      signOut: state.signOut,
      hydrate: state.hydrate,
    }),
    shallow
  );
  
  return {
    token,
    refreshToken,
    status,
    signIn,
    signOut,
    hydrate,
  };
};

export const useAuthStoreAxiosState = () => {
  const { token, refreshToken } = useAuthStore.getState();

  return { token, refreshToken };
};
