import { createClient, type RedisClientType } from "@redis/client";
import { logger } from "../utils/logConfig.js";

export class RedisMessenger {
  private readonly _subscriber: RedisClientType;
  private readonly _publisher: RedisClientType;

  constructor(
    url: string = process.env["REDIS_URI"] ?? "redis://localhost:6379",
  ) {
    this._subscriber = createClient({ url });
    this._publisher = createClient({ url });

    this._subscriber.on("error", (err) =>
      console.error("redis subscriber error", err),
    );
    this._publisher.on("error", (err) =>
      console.error("redis publisher error", err),
    );
  }

  async connect(): Promise<void> {
    await Promise.all([this._subscriber.connect(), this._publisher.connect()]);
  }

  async disconnect(): Promise<void> {
    await Promise.all([this._subscriber.quit(), this._publisher.quit()]);
  }

  async publish<T>(topic: string, message: T): Promise<number> {
    return this._publisher.publish(topic, JSON.stringify(message));
  }

  async subscribe<T>(
    topic: string,
    handler: (message: T) => void,
  ): Promise<void> {
    await this._subscriber.subscribe(topic, (message: string) => {
      let parsed: unknown;
      try {
        parsed = JSON.parse(message);

        if (typeof parsed !== "object" || parsed === null) {
          console.warn(`malformed message on ${topic}`);
          return;
        }

        handler(parsed as T);
      } catch {
        logger.warn(`malformed message on ${topic}`);
        return;
      }
    });
  }

  async unsubscribe(topic: string): Promise<void> {
    await this._subscriber.unsubscribe(topic);
  }
}
