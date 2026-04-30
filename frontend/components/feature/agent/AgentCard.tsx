import { View, XStack, YStack, styled, Text } from 'tamagui';
import { Server, ChevronRight } from '@tamagui/lucide-icons-2';
import { Badge } from '@/components/ui/Badge';
import { Agent } from '@/lib/api/models/agent';
import React from 'react';
import { router } from 'expo-router';
import { useAgentStore } from '@/store/agent/useAgentStore';

const Container = styled(YStack, {
  bg: '$color2',
  p: '$4',
  rounded: '$4',
  items: 'center',
  gap: '$3',
  borderWidth: 1,
  borderColor: '#2a2a2a',
  pressStyle: {
    scale: 0.98,
    bg: '$color3',
  },
});

const IconContainer = styled(View, {
  bg: '$color4',
  rounded: '$3',
  items: 'center',
  justify: 'center',
  width: 48,
  height: 48,
});

/*
type MetricsAction = { type: 'ADD_METRIC'; metric: AgentCpuMetric } | { type: 'CLEAR_METRICS' };

function metricsReducer(state: AgentCpuMetric[], action: MetricsAction): AgentCpuMetric[] {
  switch (action.type) {
    case 'ADD_METRIC':
      console.log('reducer', action.type, 'prev length', state.length);
      return [...state, action.metric].splice(-120);
    case 'CLEAR_METRICS':
      return [];
    default:
      return state;
  }
}
const [data, dispatch] = useReducer(metricsReducer, []);
useEffect(() => {
  if (state !== 'ACTIVE') {
    dispatch({ type: 'CLEAR_METRICS' });
    return;
  }

  const socket = getSocket();

  const handleCpuMetric = (agentId: string, metric: AgentCpuMetric) => {
    if (agentId !== identifier) return;
    dispatch({ type: 'ADD_METRIC', metric });
  };

  socket?.on('agent:metrics:cpu', handleCpuMetric);

  return () => {
    socket?.off('agent:metrics:cpu', handleCpuMetric);
    dispatch({ type: 'CLEAR_METRICS' });
  };
}, [state, identifier]);
*/
export interface AgentCardProps {
  agent: Agent;
}

export function AgentCard({ agent }: AgentCardProps) {
  const { setAgent } = useAgentStore()
  
  if (!agent) {
    return null;
  }

  return (
    <Container
      onPress={() => {
        setAgent(agent);
        router.push({
          pathname: '/agent/[id]',
          params: {
            id: agent.id,
          },
        })
      }}
    >
      <XStack gap={12} items="center">
        <IconContainer>
          <Server size={24} color="$color7" />
        </IconContainer>
        <YStack gap={4} justify="center" flex={1}>
          <Text fontWeight="bold" fontSize={18}>
            {agent.display_name}
          </Text>
          <Text color="$gray11">{agent.registered_ip}</Text>
        </YStack>

        {(() => {
          switch (agent.status) {
            case 'ACTIVE':
              return <Badge bg="$green9">Active</Badge>;
            case 'INACTIVE':
              return <Badge bg="$red9">Inactive</Badge>;
            case 'MAINTENANCE':
              return <Badge bg="$yellow9">Maintenance</Badge>;
            case 'ENROLLING':
              return <Badge bg="$blue9">Enrolling</Badge>;
            default:
              return <Badge bg="$gray9">Desconocido</Badge>;
          }
        })()}
        <ChevronRight size={18} color="$color8" />
      </XStack>

      {/*state === 'ACTIVE' && (
        <>
          <View height={100} bg="$color4" width="100%" rounded="$4" pt="$3">
            <CartesianChart
              data={data}
              xKey="timestamp"
              yKeys={['value']}
              domain={{ x: [Date.now() - 60_000, Date.now()], y: [0, 100] }}
              xAxis={{
                lineWidth: 0,
              }}
              yAxis={[
                {
                  lineWidth: 0,
                },
              ]}>
              {({ points }) => {
                return (
                  <Line
                    points={points.value}
                    color="#4786e6"
                    strokeWidth={2}
                    curveType="monotoneX"
                    animate={{ type: 'timing', duration: 100 }}
                  />
                );
              }}
            </CartesianChart>
            <Badge
              backdropFilter="blur(10px)"
              bg="rgba(180, 180, 180, 0.2)"
              position="absolute"
              b={0.15}
              l={0.16}
              opacity={0.75}
              m="$2"
              borderWidth={0}
              fontSize={12}>
              CPU
            </Badge>
          </View>
        </>
      ) */}
    </Container>
  );
}

export default React.memo(AgentCard);
