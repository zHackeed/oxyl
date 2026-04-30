import { Socket as ActiveSocketConn } from "socket.io";

export abstract class Middleware {

  constructor(private readonly _eventName: string) {}

  abstract validate(userConnection: ActiveSocketConn, ...args: any[]): Promise<void | Error>;

  get eventName(): string {
    return this._eventName;
  }

  register(userConnection: ActiveSocketConn): void {
    userConnection.use(async ([event, ...args], next) => {
      if (event !== this._eventName) return next();
      const result = await this.validate(userConnection, ...args);

      if (result instanceof Error) {
        console.log("middleware error", result);
        return next(result);
      }
      next();
    });
  }
}
