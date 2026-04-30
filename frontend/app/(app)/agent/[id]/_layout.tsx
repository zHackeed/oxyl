import { agentService } from '@/lib/service/agent';
import { useAgentFarcade } from '@/store/agent/useAgentFarcade';
import { useQuery } from '@tanstack/react-query';
import { Stack, useLocalSearchParams } from 'expo-router';
import { useEffect } from 'react';
import { useWebsocketFarcade } from '@/store/websocket/useWebsocketFarcade';

export default function AgentLayout() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { agent, status, setAgent } = useAgentFarcade();
  const { join, leave, connected } = useWebsocketFarcade();

  const { data, isLoading } = useQuery({
    queryKey: ['agent', id],
    queryFn: () => agentService.getOne(id || agent?.id || ''),
    initialData: agent,
  });

  useEffect(() => {
    if (data) {
      setAgent(data);
    }

    return () => {
      setAgent(null);
    };
  }, [data]);

  useEffect(() => {
    if (!agent?.id) return;
    if (!connected) return;

    join('agent', agent.id);

    return () => {
      if (agent?.id) {
        leave('agent', agent.id);
      }
    };
  }, [connected, agent?.id]);

  if (isLoading) {
    return null;
  }

  return (
    <Stack screenOptions={{ headerShown: false }}>
      <Stack.Protected guard={status !== 'ENROLLING'}>
        <Stack.Screen name="(tabs)" />
        <Stack.Screen
          name="info"
          options={{
            presentation: 'modal',
          }}
        />
      </Stack.Protected>
      <Stack.Protected guard={status === 'ENROLLING'}>
        <Stack.Screen name="enrolling" />
      </Stack.Protected>
    </Stack>
  );
}
