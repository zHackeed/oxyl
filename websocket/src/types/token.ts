import type { JWTPayload } from 'jose';

export type JWTTokenType = 'agent' | 'user';

export interface TokenClaims extends JWTPayload {
  identifier: string;
  holder?: string;
  type: JWTTokenType;
}