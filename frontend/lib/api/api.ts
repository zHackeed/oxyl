import { AuthStatus, useAuthStore } from '@/store/auth/useAuthStore';
import { TokenService } from '../service/token-service';
import axios, { InternalAxiosRequestConfig } from 'axios';

const API_CONFIG = {
  baseURL: "http://localhost:19999",
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
};

const caller = axios.create(API_CONFIG);

caller.interceptors.request.use(async (config: InternalAxiosRequestConfig) => {
  console.log(process.env.OXYL_API_ENDPOINT)
  const token = await TokenService.getAccessToken();
  if (token && new Date(token.expires_at) > new Date()) {
      config.headers['Authorization'] = `Bearer ${token.token}`;
  }
  return config;
});

caller.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      const refreshToken = await TokenService.getRefreshToken();
      if (!refreshToken || refreshToken.expires_at < new Date()) {
        useAuthStore.setState({
          token: null,
          refreshToken: null,
          status: AuthStatus.UNAUTHENTICATED,
        });
        return Promise.reject(error);
      }

      try {
        const response = await axios.post(`${API_CONFIG.baseURL}/auth/refresh`, {
          refresh_token: refreshToken.token,
        });

        if (response.status !== 200) {
          await TokenService.clearTokens();
          useAuthStore.setState({
            token: null,
            refreshToken: null,
            status: AuthStatus.UNAUTHENTICATED,
          });
          return Promise.reject(new Error('Failed to refresh token'));
        }

        // Update tokens in storage
        await TokenService.setAccessToken(response.data.access_token);
        await TokenService.setRefreshToken(response.data.refresh_token);

        useAuthStore.setState({
          token: response.data.access_token,
          refreshToken: response.data.refresh_token,
          status: AuthStatus.AUTHENTICATED,
        });
      } catch (error) {
        console.error('Failed to refresh token:', error);
        useAuthStore.setState({
          token: null,
          refreshToken: null,
          status: AuthStatus.UNAUTHENTICATED,
        });
        return Promise.reject(error);
      }

      // retry the original request with the new access token that might be renewed?
      return caller(originalRequest);
    }

    return Promise.reject(error);
  }
);

export default caller;
