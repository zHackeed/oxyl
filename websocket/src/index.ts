import { Server } from "socket.io";
import { logger } from "./utils/logConfig.js";
import { TokenService } from "./service/tokenService.js";
import { AuthMiddleware } from "./middleware/authMiddleware.js";
import { RedisMessenger } from "./db/redisConn.js";
import type { CompanyUpdateActions, SocketMetadata, UserSocketReq } from "./types/socket.js";
import { CompanyMiddleware } from "./middleware/companyMiddleware.js";
import { CompanyService } from "./service/companyService.js";
import { UserInvalidationService } from "./service/invalidationService.js";

const io = new Server<
  UserSocketReq,
  CompanyUpdateActions,
  SocketMetadata
>({
  path: "/ws",
  cors: {
    origin: "*",
  },
  connectionStateRecovery: {
    maxDisconnectionDuration: 2000,
    skipMiddlewares: false,
  },
  cleanupEmptyChildNamespaces: true,
});

const redisMessenger = new RedisMessenger(process.env["REDIS_URL"]);
await redisMessenger.connect();

new UserInvalidationService(io, redisMessenger);

const companyMiddleware = new CompanyMiddleware('join');
const tokenService = await TokenService.create("/data/keys");
const authMiddleware = new AuthMiddleware(tokenService);
const companyService = new CompanyService(redisMessenger, io);


io.use(authMiddleware.handle);
io.on("connection", (socket) => {
  logger.info("User connected", socket.id);

  const userId = socket.data.userId;
  socket.join(`user:${userId}`);

  companyMiddleware.register(socket);

  socket.on("join", (type, id) => {
    logger.info("user has requested to join", type, id);
    switch (type) {
      case "company":
        companyService.addListener(id, socket);
        break;
      case "agent":
        //todo: agent logic
        break;
    }
  });

  socket.on("leave", (type, id) => {
    logger.info("user has requested to leave", type, id);
    switch (type) {
      case "company":
        companyService.removeListener(id, socket);
        break;
      case "agent":
        //todo: agent logic
        break;
    }
  });

  
  socket.on("disconnecting", (reason) => {
    logger.info("User is disconnecting", socket.id, reason);
    socket.rooms.forEach((room) => {
      if (room.startsWith("company:")) {
        companyService.removeListener(room.replace("company:", ""), socket);
      }
    });
  });


  socket.on("disconnect", () => {
    socket.leave(`user:${userId}`);
    logger.info("User disconnected", socket.id, socket.conn.remoteAddress);
  });
});

// thonking all the time
io.listen(19977);
