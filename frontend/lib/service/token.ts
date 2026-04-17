import * as SecureStorage from 'expo-secure-store';
import { AuthToken } from '@/lib/api/models/token';

export const TokenService = {
  async setAccessToken(accessTokenData: AuthToken) {
    await SecureStorage.setItemAsync('oxylAccessToken', JSON.stringify(accessTokenData));
  },

  async getAccessToken(): Promise<AuthToken | null> {
    const token = await SecureStorage.getItemAsync('oxylAccessToken');

    if (!token) {
      return null;
    }

    const parsed = JSON.parse(token);

    return {
      ...parsed,
      expires_at: new Date(parsed.expires_at),
    };
  },

  async setRefreshToken(refreshTokenData: AuthToken) {
    await SecureStorage.setItemAsync('oxylRefreshToken', JSON.stringify(refreshTokenData));
  },

  async getRefreshToken(): Promise<AuthToken | null> {
    const token = await SecureStorage.getItemAsync('oxylRefreshToken');

    if (!token) {
      return null;
    }

    const parsed = JSON.parse(token);

    return {
      ...parsed,
      expires_at: new Date(parsed.expires_at),
    };
  },

  async clearTokens() {
    await SecureStorage.deleteItemAsync('oxylAccessToken');
    await SecureStorage.deleteItemAsync('oxylRefreshToken');
  },
};
