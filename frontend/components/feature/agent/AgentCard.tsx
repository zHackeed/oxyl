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
    </Container>
  );
}

export default React.memo(AgentCard);
