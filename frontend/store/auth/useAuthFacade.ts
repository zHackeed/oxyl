import useAuthStore from './useAuthStore';
import { shallow } from 'zustand/shallow';

// We need this to expose the methods from zustard
export const useAuthFacade = () => {
  const { token, refreshToken, status, signIn, signOut } = useAuthStore(
    (state) => ({
      token: state.token,
      refreshToken: state.refreshToken,
      status: state.status,
      signIn: state.signIn,
      signOut: state.signOut,
    }),
    shallow
  );

  return {
    token,
    refreshToken,
    status,
    signIn,
    signOut,
  };
};