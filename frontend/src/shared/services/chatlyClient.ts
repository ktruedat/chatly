import { SERVER_URL } from "@/config/api.config";
import type { ClientMessage, ServerMessage } from "@/shared/types/ws";

type EventHandler = (payload?: any) => void;

export class ChatlyClient {
  private ws?: WebSocket;
  private url: string;
  private shouldReconnect = true;
  private reconnectAttempts = 0;
  private reconnectTimer: number | null = null;
  private readonly maxReconnectAttempts = 10;
  private readonly baseDelayMs = 500;
  private readonly maxDelayMs = 30_000;
  private listeners: {
    open: EventHandler[];
    close: EventHandler[];
    error: EventHandler[];
    message: ((msg: ServerMessage) => void)[];
  } = { open: [], close: [], error: [], message: [] };

  constructor(url?: string) {
    this.url = this.buildWsUrl(url || SERVER_URL);
  }

  private buildWsUrl(provided?: string): string {
    if (provided && provided.startsWith("http")) {
      try {
        const u = new URL(provided);
        const proto = u.protocol === "https:" ? "wss:" : "ws:";
        const base = `${proto}//${u.host}`;
        return `${base}/ws`;
      } catch (e) {
        console.error("Invalid WS URL provided:", e);
      }
    }
    if (
      provided &&
      (provided.startsWith("ws:") || provided.startsWith("wss:"))
    ) {
      return provided;
    }
    if (typeof window !== "undefined") {
      const proto = window.location.protocol === "https:" ? "wss:" : "ws:";
      return `${proto}//${window.location.host}/ws`;
    }
    return "ws://localhost:3000/ws";
  }

  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) return resolve();
      if (this.reconnectTimer) {
        clearTimeout(this.reconnectTimer);
        this.reconnectTimer = null;
      }

      this.ws = new WebSocket(this.url);

      this.ws.onopen = () => {
        this.emitOpen();
        resolve();
      };

      this.ws.onclose = (ev) => {
        this.emitClose(ev);
        if (this.shouldReconnect) this.scheduleReconnect();
      };

      this.ws.onerror = (ev) => {
        this.emitError(ev);
        reject(new Error("WebSocket connection error"));
      };

      this.ws.onmessage = (ev) => {
        try {
          const parsed = JSON.parse(ev.data) as ServerMessage;
          this.emitMessage(parsed);
        } catch (e) {
          console.error("Invalid WS JSON:", e);
        }
      };
    });
  }

  disconnect() {
    this.shouldReconnect = false;
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    if (this.ws) {
      try {
        this.ws.close();
      } catch (e) {
        console.log("WebSocket close error:", e);
      }
      this.ws = undefined;
    }
  }

  isConnected() {
    return !!(this.ws && this.ws.readyState === WebSocket.OPEN);
  }

  sendMessage(msg: ClientMessage) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      throw new Error("WebSocket is not connected");
    }
    this.ws.send(JSON.stringify(msg));
  }

  private scheduleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) return;
    this.reconnectAttempts += 1;
    const exp = Math.min(
      this.baseDelayMs * 2 ** (this.reconnectAttempts - 1),
      this.maxDelayMs
    );
    const jitter = Math.floor(exp * (Math.random() * 0.4 - 0.2));
    const delay = Math.max(100, exp + jitter);
    this.reconnectTimer = window.setTimeout(() => {
      this.reconnectTimer = null;
      this.connect().catch(() => {
        this.scheduleReconnect();
      });
    }, delay);
  }

  onOpen(cb: EventHandler) {
    this.listeners.open.push(cb);
  }

  onClose(cb: EventHandler) {
    this.listeners.close.push(cb);
  }

  onError(cb: EventHandler) {
    this.listeners.error.push(cb);
  }

  onMessage(cb: (msg: ServerMessage) => void) {
    this.listeners.message.push(cb);
  }

  private emitOpen() {
    this.listeners.open.forEach((cb) => cb());
  }

  private emitClose(ev?: any) {
    this.listeners.close.forEach((cb) => cb(ev));
  }

  private emitError(ev?: any) {
    this.listeners.error.forEach((cb) => cb(ev));
  }

  private emitMessage(msg: ServerMessage) {
    this.listeners.message.forEach((cb) => cb(msg));
  }
}

export default ChatlyClient;
