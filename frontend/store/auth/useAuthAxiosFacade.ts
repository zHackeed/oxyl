import useAuthStore from "./useAuthStore";


// TODO: Fix require cycle
export const useAuthStoreAxiosState = () => {
  const { token, refreshToken, refreshTokens, signOut } = useAuthStore.getState();

  return { token, refreshToken, refreshTokens, signOut };
};
