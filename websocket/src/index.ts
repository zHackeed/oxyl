import { Server } from "socket.io";
import { logger } from "./utils/logConfig.js";
import { TokenService } from "./service/tokenService.js";
import { AuthMiddleware } from "./middleware/authMiddleware.js";
import { RedisMessenger } from "./db/redisConn.js";


const io = new Server({
  path: "/ws",
  cors: {
    origin: "*",
  },
  connectionStateRecovery: {
    maxDisconnectionDuration: 2000,
    skipMiddlewares: false,
  },
  cleanupEmptyChildNamespaces: true
})

const redisMessenger = new RedisMessenger(process.env.REDIS_URL);
const tokenService = await TokenService.create("/data/keys");
const authMiddleware = new AuthMiddleware(tokenService);

io.use(authMiddleware.handle)
io.on("connection", (socket) => {
  logger.info("User connected", socket.id);
  
  const userId = socket.data.userId;
  socket.join(`user-connections:${userId}`)
  
  // * todo: disconnection logic and handling of leaving a room. Might have to think if we have to use a service for state tracking?
  
  socket.on("disconnect", () => {
    socket.leave(`user-connections:${userId}`);
    logger.info("User disconnected", socket.id);
  });
});


// thonking all the time
io.listen(19977);