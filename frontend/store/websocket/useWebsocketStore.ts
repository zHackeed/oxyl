
import { Socket, io } from "socket.io-client";
import { createWithEqualityFn } from "zustand/traditional";
import useAuthStore from "../auth/useAuthStore";
import { CompanyUpdateActions, RoomType, UserSocketReq } from "@/lib/websocket/actions";

let _socket: Socket<CompanyUpdateActions, UserSocketReq> | null = null;
const _subscriptions = new Set<string>();

export function getSocket() {
  return _socket;
}

export interface WebsocketStateProps {
  connected: boolean;
  connect: () => void;
  join: (type: RoomType, id: string) => void;
  leave: (type: RoomType, id: string) => void;
  disconnect: () => void;
}

export const useWebsocketStore = createWithEqualityFn<WebsocketStateProps>()((set) => ({
  connected: false,
  connect: () => {
    const token = useAuthStore.getState().token;
    if (!token) return;

    _socket = io("http://10.0.60.4:19977", {
      auth: (cb) => cb({ token: useAuthStore.getState().token?.token || "" }),
      path: "/ws",
      autoConnect: true,
      reconnectionDelayMax: 10000,
      reconnectionAttempts: 10,
    });

    _socket.on("connect", () => {
      set({ connected: true })

      if (_socket?.recovered) return;

      for (const value of _subscriptions) {
        const [type, id] = value.split("|") as [RoomType, string];

        _socket?.emit("join", type, id);
      }
    });
    _socket.on("connect_error", (error) => {
      console.log(error)

      if (error.message !== "Unauthorized") return;
      //As the socket is not on the proper state for reconnection, it would never try to reconnect automatically
      //So we need to manually trigger the reconnection.
      setTimeout(() => {
        _socket?.connect();
      }, 2000);
    });
    _socket.on("disconnect", () => {
      set({ connected: false }) 

      if (!_socket?.active) {
        setTimeout(() => {
          _socket?.connect()

          set({
            connected: true
          })
        }, 2000)
      }
    });
  },

  join: (type: RoomType, id: string) => {
    _subscriptions.add(`${type}|${id}`)
    console.log("joining", type, id)
    _socket?.emit("join", type, id);
  },

  leave: (type: RoomType, id: string) => {
    _subscriptions.delete(`${type}|${id}`)
    _socket?.emit("leave", type, id);
  },

  disconnect: () => {
    _socket?.disconnect();
    _socket = null;
    set({ connected: false });
  },
}));
