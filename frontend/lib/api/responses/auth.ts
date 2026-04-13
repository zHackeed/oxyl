
import { AuthToken } from "@/lib/api/models/token";

type AuthTokenResponse = {
  access_token: AuthToken;
  refresh_token: AuthToken;
};

export { AuthTokenResponse };
