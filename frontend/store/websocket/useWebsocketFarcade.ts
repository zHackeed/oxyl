import { useWebsocketStore } from './useWebsocketStore';
import { shallow } from 'zustand/shallow';

export const useWebsocketFarcade = () => {
  const { join, leave, connected } = useWebsocketStore(
    (state) => ({
      join: state.join,
      leave: state.leave,
      connected: state.connected,
    }),
    shallow
  );
  
  return {
    join,
    leave,
    connected,
  };
};