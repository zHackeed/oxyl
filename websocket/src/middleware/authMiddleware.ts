import type { ExtendedError, Socket } from "socket.io";
import type { TokenService } from "../service/tokenService.js";

export class AuthMiddleware {
  constructor(private readonly tokenService: TokenService) { }
  
  handle = async (socket: Socket, next: (err?: ExtendedError) => void) => {
    const token =
      (socket.handshake.auth?.token as string | undefined) ??
      socket.handshake.headers.authorization;
    
    if (!token) {
      return next(new Error("Unauthorized: No token provided"));
    }
    
    try {
      const claims = await this.tokenService.verifyToken(token);
      
      socket.data.userId = claims.identifier
      socket.data.jti = claims.jti
      
      next();
    } catch (e) {
      next(new Error("Unauthorized"));
    }
  }
}