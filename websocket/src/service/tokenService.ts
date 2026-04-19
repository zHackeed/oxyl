import type { TokenClaims } from "../types/token.js";
import { readFileSync } from "node:fs";
import * as jose from "jose";
import { InvalidTokenError } from "../types/error.js";
import { logger } from "../utils/logConfig.js";

const TOKEN_ISSUER = "oxyl";
const ALLOWED_AUDIENCES = [
  "https://api.oxyl.zhacked.me",
  "https://ingress.oxyl.zhacked.me"
];

export class TokenService {
  private constructor(private readonly publicKey: jose.CryptoKey) {}

  static async create(keyLocs: string): Promise<TokenService> {
    const publicKey = readFileSync(`${keyLocs}/ed25519-pub.pem`);
    const key = await jose.importSPKI(publicKey.toString(), "EdDSA");
    return new TokenService(key);
  }

  async verifyToken(token: string): Promise<TokenClaims> {
    const sanitized = token.trim().replace(/^Bearer\s+/i, '');
    if (!sanitized) throw new InvalidTokenError('empty token');
    
    try {
      const { payload } = await jose.jwtVerify(sanitized, this.publicKey, {
        algorithms: ["EdDSA"],
        issuer: TOKEN_ISSUER,
        audience: ALLOWED_AUDIENCES,
      });
    
      const claims = payload as TokenClaims;
      
      logger.info(`verified token for user ${claims.type}`);

      if (claims.type != 'USER') {
        throw new InvalidTokenError('invalid token for this service');
      }

      return claims;
    } catch (e) {
      throw new InvalidTokenError(e instanceof Error ? e.message : 'failed to verify token');
    }
  }
}