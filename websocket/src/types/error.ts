
export class InvalidTokenError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "InvalidTokenError";
  }
}

export class PermissionDeniedError extends Error {
  constructor() {
    super('user does not have the required permissions to perform this action');
    this.name = 'PermissionDeniedError';
  }
}
