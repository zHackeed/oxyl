import { AuthToken } from '@/lib/api/models/token';
import { AuthService } from '@/lib/service/auth';
import { TokenService } from '@/lib/service/token';
import { createWithEqualityFn } from 'zustand/traditional';

export enum AuthStatus {
  AUTHENTICATED = 'authenticated',
  LOADING = 'loading',
  UNAUTHENTICATED = 'unauthenticated',
}

export interface AuthStatProps {
  token: AuthToken | null;
  refreshToken: AuthToken | null;
  status: AuthStatus;

  signIn: (email: string, password: string) => Promise<boolean>;
  signOut: () => void;
  refreshTokens: (token: AuthToken, refreshToken: AuthToken) => Promise<void>;
}

const initialAuthState = {
  token: null,
  refreshToken: null,
  status: AuthStatus.LOADING,
};

const useAuthStore = createWithEqualityFn<AuthStatProps>()((set) => {
  const initialState = initialAuthState;

  const getAccessToken = async () => {
    const token = await TokenService.getAccessToken();
    return token;
  };

  const getRefreshToken = async () => {
    const token = await TokenService.getRefreshToken();
    return token;
  };

  Promise.all([getRefreshToken(), getAccessToken()])
    .then(([refreshToken, token]) => {
      if (refreshToken && token) {
        set((state) => ({
          ...state,
          token: token,
          refreshToken: refreshToken,
          status: AuthStatus.AUTHENTICATED,
        }));
      } else {
        set((state) => ({
          ...state,
          status: AuthStatus.UNAUTHENTICATED,
        }));
      }
    })
    .catch((error) => {
      set((state) => ({
        ...state,
        status: AuthStatus.UNAUTHENTICATED,
      }));
    });

  return {
    ...initialState,
    signIn: async (email: string, password: string) => {
      const result = await AuthService.login(email, password); // to do, remove
      if (result) {
        set((state) => ({
          ...state,
          token: result.access_token,
          refreshToken: result.refresh_token,
          status: AuthStatus.AUTHENTICATED,
        }));
      }

      return false;
    },

    refreshTokens: async (token, refreshToken) => {
      await TokenService.setAccessToken(token);
      await TokenService.setRefreshToken(refreshToken);
      set((state) => ({
        ...state,
        token: token,
        refreshToken: refreshToken,
      }));
    },

    signOut: async () => {
      await TokenService.clearTokens();
      set((state) => ({
        ...state,
        token: null,
        refreshToken: null,
        status: AuthStatus.UNAUTHENTICATED,
      }));
    },
  };
});

export default useAuthStore;
