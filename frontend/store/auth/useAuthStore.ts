import { AuthToken } from '@/lib/api/models/token';
import { AuthService } from '@/lib/service/auth-service';
import { TokenService } from '@/lib/service/token-service';
import { createWithEqualityFn } from 'zustand/traditional'

export enum AuthStatus {
  AUTHENTICATED = 'authenticated',
  LOADING = 'loading',
  UNAUTHENTICATED = 'unauthenticated',
}

export interface AuthState {
  token: AuthToken | null;
  refreshToken: AuthToken | null;
  status: AuthStatus;
  signIn: (email: string, password: string) => Promise<void>;
  signOut: () => Promise<void>;
  hydrate: () => Promise<boolean>;
}

export const useAuthStore = createWithEqualityFn<AuthState>((set) => ({
  token: null,
  refreshToken: null,
  status: AuthStatus.LOADING,
  signIn: async (email: string, password: string) => {
    const result = await AuthService.login(email, password);
    if (result) {
      set({
        token: result.access_token,
        refreshToken: result.refresh_token,
        status: AuthStatus.AUTHENTICATED,
      });
    }
  },
  signOut: async () => {
    await AuthService.logout();
    set({
      token: null,
      refreshToken: null,
      status: AuthStatus.UNAUTHENTICATED,
    });
  },
  hydrate: async () => {
    try {
      const accessToken = await TokenService.getAccessToken();
      const refreshToken = await TokenService.getRefreshToken();

      if (accessToken && refreshToken) {
        if (refreshToken.expires_at < new Date()) {
          // refresh token is expired, we can't do anything
          return true;
        }

        // are the access token expired?
        if (accessToken.expires_at < new Date()) {
          // refresh the token
          const result = await AuthService.refresh(refreshToken.token);
          if (result) {
            set({
              token: result.access_token,
              refreshToken: result.refresh_token,
              status: AuthStatus.AUTHENTICATED,
            });
          }

          return true;
        }
        set({
          token: accessToken,
          refreshToken: refreshToken,
          status: AuthStatus.AUTHENTICATED,
        });
        return true;
      } else {
        set({
          status: AuthStatus.UNAUTHENTICATED,
        });
        return true;
      }
    } catch (error) {
      console.error('Error hydrating auth state:', error);
      return false;
    }
  },
}));


