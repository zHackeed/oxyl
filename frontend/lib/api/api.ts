import { TokenService } from '../service/token';
import axios, { InternalAxiosRequestConfig } from 'axios';
import { useAuthStoreAxiosState } from '@/store/auth/useAuthAxiosFacade';

const AUTH_BASE_URL = 'http://10.0.60.5:19999/auth';

const API_CONFIG = {
  // Todo: make configurable
  baseURL: 'http://10.0.60.5:19999/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
};

const Caller = axios.create(API_CONFIG);

const AuthCaller = axios.create({
  baseURL: 'http://10.0.60.5:19999/',
  withCredentials: false,
});

Caller.interceptors.request.use(async (config: InternalAxiosRequestConfig) => {
  const { token } = useAuthStoreAxiosState();
  if (token && new Date(token.expires_at) > new Date()) {
    config.headers['Authorization'] = `Bearer ${token.token}`;
  }

  return config;
});

let refreshPromise: Promise<void> | null = null;

Caller.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      const { refreshTokens, signOut, refreshToken } = useAuthStoreAxiosState();

      if (!refreshToken || refreshToken.expires_at < new Date()) {
        await signOut();
        return Promise.reject(error);
      }

      if (!refreshPromise) {
        refreshPromise = axios
          .post(`${AUTH_BASE_URL}/refresh`, {
            refresh_token: refreshToken.token,
          })
          .then(async (response) => {
            console.log('refresing');
            if (response.status !== 200) {
              await TokenService.clearTokens();
              signOut();
              return Promise.reject(new Error('Failed to refresh token'));
            }
            // Update tokens in storage
            await refreshTokens(response.data.access_token, response.data.refresh_token);
          })
          .catch(async (error) => {
            console.error('Failed to refresh token:', error);
            await signOut();
            return Promise.reject(error);
          })
          .finally(() => {
            refreshPromise = null;
          });
      }

      try {
        await refreshPromise;
      } catch (error) {
        return Promise.reject(error);
      }

      // retry the original request with the new access token that might be renewed?
      return Caller(originalRequest);
    }

    return Promise.reject(error);
  }
);

export { Caller, AuthCaller };
