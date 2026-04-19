import type { JWTPayload } from 'jose';

export type JWTTokenType = 'AGENT' | 'USER';

export interface TokenClaims extends JWTPayload {
  identifier: string;
  holder?: string;
  type: JWTTokenType;
}