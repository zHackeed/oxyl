import caller from '@/lib/api/api';
import { LogoutRequest, UserLoginRequest, UserRegisterRequest } from '@/lib/api/requests/user';
import { AuthTokenResponse } from '../api/responses/auth';
import { TokenService } from './token-service';

export const AuthService = {
  async login(email: string, password: string): Promise<AuthTokenResponse | null> {
    try {
      const response = await caller.post('/auth/login', {
        email,
        password,
      } as UserLoginRequest);


      if (response.status !== 200) {
        // Check params
        return null;
      }

      const authToken = response.data as AuthTokenResponse;

      await TokenService.setAccessToken(authToken.access_token);
      await TokenService.setRefreshToken(authToken.refresh_token);
      return authToken;
    } catch (error) {
      console.error(error);
      return null;
    }
  },

  // The user when it gets registered it returns a 201 created with the user.
  // If it was successfull, we return the user to the login page after a few secs with a banner that would not be hidden
  async register(name: string, surname: string, email: string, password: string): Promise<void> {
    try {
      const response = await caller.post('/auth/register', {
        name,
        surname,
        email,
        password,
      } as UserRegisterRequest);

      if (response.status !== 201) {
        throw new Error('Registration failed');
      }

      // might force the user to call login?
    } catch (error) {
      console.error(error);
      throw error;
    }
  },

  async refresh(refreshToken: string): Promise<AuthTokenResponse | null> {
    try {
      const response = await caller.post('/auth/refresh', { refresh_token: refreshToken });

      if (response.status !== 200) {
        //todo, maybe throw an error?
        return null;
      }

      const authToken = response.data as AuthTokenResponse;

      await TokenService.setAccessToken(authToken.access_token);
      await TokenService.setRefreshToken(authToken.refresh_token);

      return authToken;
    } catch (error) {
      console.error(error);
      throw error;
    }
  },

  async logout(): Promise<void> {
    try {
      const refreshToken = await TokenService.getRefreshToken();
      if (!refreshToken) {
        // No refresh token, just clear all tokens
        await TokenService.clearTokens();
        return;
      }

      if (refreshToken.expires_at < new Date()) {
        // Token expired, clear it and return
        await TokenService.clearTokens();
        return;
      }

      const response = await caller.post('/auth/logout', {
        refresh_token: refreshToken.token,
      } as LogoutRequest);

      if (response.status !== 200) {
        throw new Error('Logout failed');
      }

      await TokenService.clearTokens();
    } catch (error) {
      console.error(error);
    }
  },
};
