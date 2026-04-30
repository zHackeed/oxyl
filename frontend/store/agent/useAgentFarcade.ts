import { useAgentStore } from "./useAgentStore";
import { shallow } from "zustand/shallow";

export const useAgentFarcade = () => {
  const { agent, setAgent, status, setStatus } = useAgentStore(
    (state) => ({
      agent: state.agent,
      setAgent: state.setAgent,
      status: state.status,
      setStatus: state.setStatus,
    }),
    shallow
  );
  
  return {
    agent,
    setAgent,
    status,
    setStatus,
  };
};