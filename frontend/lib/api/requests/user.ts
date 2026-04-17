type UserLoginRequest = {
  email: string;
  password: string;
};

type UserRegisterRequest = {
  name: string;
  surname: string;
  email: string;
  password: string;
  confirmPassword: string;
};

type RefreshTokenRequest = {
  refresh_token: string;
};

type LogoutRequest = {
  refresh_token: string;
};

export { UserLoginRequest, UserRegisterRequest, RefreshTokenRequest, LogoutRequest };
