import { Socket as ActiveSocketConn, Server } from "socket.io";

export abstract class Middleware {
  private readonly _server: Server;
  private readonly _eventName: string;

  constructor(server: Server, eventName: string) {
    this._server = server;
    this._eventName = eventName;
  }

  abstract validate(userConnection: ActiveSocketConn): void | Error;

  get eventName(): string {
    return this._eventName;
  }

  get server(): Server {
    return this._server;
  }

  register(): void {
    this._server.on('connection', (userConnection) => {
      userConnection.use(([event, ...args], next) => {
        if (event !== this._eventName) return next();
        const result = this.validate(userConnection);
      
        if (result instanceof Error) return next(result);  
        next();
      });
    });
  }
}
