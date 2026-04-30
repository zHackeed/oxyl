export class Cache<K, V> {
  private cache = new Map<K, V>();
  private timers = new Map<K, ReturnType<typeof setTimeout>>();

  constructor(private maxSize: number = 1000) {}

  get(key: K): V | undefined {
    return this.cache.get(key);
  }

  set(key: K, value: V, ttl?: number): void {
    const existing = this.timers.get(key);
    if (existing) clearTimeout(existing);

    if (!this.cache.has(key) && this.cache.size >= this.maxSize) {
      const firstKey = this.cache.keys().next().value;
      if (firstKey !== undefined) this.evict(firstKey);
    }

    this.cache.set(key, value);

    if (ttl) {
      this.timers.set(key, setTimeout(() => this.evict(key), ttl));
    }
  }

  private evict(key: K): void {
    this.cache.delete(key);
    const timer = this.timers.get(key);
    if (timer) {
      clearTimeout(timer);
      this.timers.delete(key);
    }
  }

  delete(key: K): boolean {
    this.evict(key);
    return true;
  }

  clear(): void {
    this.timers.forEach(clearTimeout);
    this.timers.clear();
    this.cache.clear();
  }

  size(): number {
    return this.cache.size;
  }
}