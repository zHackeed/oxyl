import { createClient, type RedisClientType } from "@redis/client";

export class RedisMessenger {
  private readonly _subscriber: RedisClientType;
  private readonly _publisher: RedisClientType;

  constructor(url: string = process.env.REDIS_URL ?? 'redis://localhost:6379') {
    this._subscriber = createClient({ url });
    this._publisher = createClient({ url });

    this._subscriber.on('error', (err) => console.error('redis subscriber error', err));
    this._publisher.on('error', (err) => console.error('redis publisher error', err));
  }

  async connect(): Promise<void> {
    await Promise.all([this._subscriber.connect(), this._publisher.connect()]);
  }

  async disconnect(): Promise<void> {
    await Promise.all([this._subscriber.quit(), this._publisher.quit()]);
  }

  async publish(topic: string, message: unknown): Promise<number> {
    return this._publisher.publish(topic, JSON.stringify(message));
  }

  // ? todo: might have to check this. Right now the code diferences between both platforms might not handle properly the json parsing. As right now,
  // ? we just try to do this and see if it works... lol.
  async subscribe(topic: string, handler: (message: unknown) => void): Promise<void> {
    await this._subscriber.subscribe(topic, (message: string) => {
      let parsed: unknown;
      try {
        parsed = JSON.parse(message);
      } catch {
        console.warn(`malformed message on ${topic}`);
        return;
      }
      handler(parsed);
    });
  }

  async unsubscribe(topic: string): Promise<void> {
    await this._subscriber.unsubscribe(topic);
  }
}