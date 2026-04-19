import type { ExtendedError, Socket } from "socket.io";
import type { TokenService } from "../service/tokenService.js";
import { logger } from "../utils/logConfig.js";

export class AuthMiddleware {
  constructor(private readonly tokenService: TokenService) {}

  handle = async (socket: Socket, next: (err?: ExtendedError) => void) => {
    logger.info("incoming authentication", socket.handshake)
    const token =
      (socket.handshake.auth["token"] as string | undefined) ??
      socket.handshake.headers.authorization;

    if (!token) {
      logger.warn(
        "No token provided",
        socket.handshake.headers.authorization,
        socket.client.conn.remoteAddress,
      );
      return next(new Error("Unauthorized: No token provided"));
    }

    try {
      const claims = await this.tokenService.verifyToken(token);

      socket.data.userId = claims.identifier;
      socket.data.jti = claims.jti;

      next();
    } catch (e) {
      console.log("invalid authentication from client", e, socket.client.conn.remoteAddress);
      next(new Error("Unauthorized"));
    }
  };
}
