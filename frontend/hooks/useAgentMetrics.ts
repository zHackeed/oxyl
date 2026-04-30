import {
  AgentGeneralMetric,
  AgentMetricEntry,
  AgentMountPointMetric,
  AgentNetworkMetric,
  AgentPhysicalDiskMetric,
} from '@/lib/api/models/metrics';
import { agentService } from '@/lib/service/agent';
import { useWebsocketFarcade } from '@/store/websocket/useWebsocketFarcade';
import { getSocket } from '@/store/websocket/useWebsocketStore';
import { useQuery } from '@tanstack/react-query';
import { useEffect, useReducer } from 'react';

type MetricAction<T> =
  | { type: 'SEED'; payload: T[] }
  | { type: 'APPEND'; payload: { point: T; cap: number } }
  | { type: 'CLEAR' };

type RecordAction<T> =
  | { type: 'SEED'; payload: Record<string, T[]> }
  | { type: 'APPEND'; payload: { key: string; point: T; cap: number } }
  | { type: 'CLEAR' };

function mapRecordReducer<T>(
  state: Record<string, T[]>,
  action: RecordAction<T>
): Record<string, T[]> {
  switch (action.type) {
    case 'SEED':
      return action.payload;
    case 'APPEND':
      return {
        ...state,
        [action.payload.key]: [...(state[action.payload.key] ?? []), action.payload.point].slice(
          -action.payload.cap
        ),
      };
    case 'CLEAR':
      return {};
  }
}

function listReducer<T>(state: T[], action: MetricAction<T>): T[] {
  switch (action.type) {
    case 'SEED':
      return action.payload;
    case 'APPEND':
      return [...state, action.payload.point].slice(-action.payload.cap);
    case 'CLEAR':
      return [];
  }
}

const INTERVAL_CAPS: Record<string, number> = {
  '15m': 900,
  '1h': 240,
  '6h': 360,
  '7d': 168,
};

export function useAgentMetrics(id: string, interval: string) {
  const { connected } = useWebsocketFarcade();
  const [general, dispatchGeneral] = useReducer(listReducer<AgentGeneralMetric>, []);
  const [network, dispatchNetwork] = useReducer(mapRecordReducer<AgentNetworkMetric>, {});
  const [mounts, dispatchMounts] = useReducer(mapRecordReducer<AgentMountPointMetric>, {});

  const { data, isLoading, isLoadingError } = useQuery({
    queryKey: ['agent-metrics', id, interval],
    queryFn: () => agentService.fetchAgentMetrics(id, interval),
  });

  useEffect(() => {
    if (!data) return;
    dispatchGeneral({ type: 'SEED', payload: data.general_metrics });
    dispatchNetwork({ type: 'SEED', payload: data.network_metrics });
    dispatchMounts({ type: 'SEED', payload: data.mount_point_metrics });
  }, [data]);

  useEffect(() => {
    if (!connected) return;
    const socket = getSocket();
    if (!socket) return;

    const onMessage = (message: AgentMetricEntry) => {
      const cap = INTERVAL_CAPS[interval];
      dispatchGeneral({ type: 'APPEND', payload: { point: message.general_metrics, cap } });
      for (const point of message.network_metrics) {
        dispatchNetwork({ type: 'APPEND', payload: { key: point.interface_name, point, cap } });
      }
      for (const point of message.mounted_metrics) {
        dispatchMounts({ type: 'APPEND', payload: { key: point.mount_point, point, cap } });
      }
    };

    socket.on('agent:metric:append', onMessage);

    return () => {
      dispatchGeneral({ type: 'CLEAR' });
      dispatchNetwork({ type: 'CLEAR' });
      dispatchMounts({ type: 'CLEAR' });
      socket.off('agent:metric:append', onMessage);
    };
  }, [connected]);

  return { general, network, mounts, isLoading, isLoadingError };
}
