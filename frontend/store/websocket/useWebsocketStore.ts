
import { Socket, io } from "socket.io-client";
import { createWithEqualityFn } from "zustand/traditional";
import useAuthStore from "../auth/useAuthStore";
import { CompanyUpdateActions, RoomType, UserSocketReq } from "@/lib/websocket/actions";

let _socket: Socket<CompanyUpdateActions, UserSocketReq> | null = null;

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

    _socket = io("http://127.0.0.1:19977", {
      auth: (cb) => cb({ token: useAuthStore.getState().token?.token || "" }),
      path: "/ws",
      autoConnect: true,
      reconnectionDelayMax: 10000,
      reconnectionAttempts: 5,
    });

    _socket.on("connect", () => set({ connected: true }));
    _socket.on("connect_error", (error) => {
      console.log(error)

      if (error.message !== "Unauthorized") return;
      //As the socket is not on the proper state for reconnection, it would never try to reconnect automatically
      //So we need to manually trigger the reconnection.
      setTimeout(() => {
        _socket?.connect();
      }, 2000);
    });
    _socket.on("disconnect", () => set({ connected: false }));
  },

  join: (type: RoomType, id: string) => {
    console.log("socket on join:", _socket?.id, _socket?.connected);
    _socket?.emit("join", type, id);
  },

  leave: (type: RoomType, id: string) => {
    _socket?.emit("leave", type, id);
  },

  disconnect: () => {
    _socket?.disconnect();
    _socket = null;
    set({ connected: false });
  },
}));
