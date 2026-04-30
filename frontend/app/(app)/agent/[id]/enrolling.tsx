import { WrappedView } from '@/components/ui/WrappedView';
import { AlarmSmoke, Copy } from '@tamagui/lucide-icons-2';
import { Separator, Spinner, styled, Text, XStack, YStack } from 'tamagui';
import { ScrollView } from 'react-native';
import { useAgentFarcade } from '@/store/agent/useAgentFarcade';
import * as Clipboard from 'expo-clipboard';
import { BackButton } from '@/components/ui/BackButton';
import { useEffect } from 'react';
import { useWebsocketFarcade } from '@/store/websocket/useWebsocketFarcade';
import { getSocket } from '@/store/websocket/useWebsocketStore';
import { AgentState } from '@/lib/api/models/agent';

const Info = styled(Text, {
  color: '$gray10',
  textAlign: 'center',
});

const MonoText = styled(Text, {
  style: { fontFamily: 'Courier New' },
});

const commands = (agentId: string) => [
  `curl https://cdn.zhacked.me/binaries/oxyl-agent -o /usr/local/bin/oxyl-agent`,
  `chmod +x /usr/local/bin/oxyl-agent`,
  `oxyl-agent init -i ${agentId}`,
];

// Ismael 형 me ayudó con el diseño porque le sabe y yo no :(
export default function StartEnrollment() {
  const { agent, status, setAgent, setStatus } = useAgentFarcade();
  const cmds = commands(agent?.id ?? '');
  const { connected } = useWebsocketFarcade();


  useEffect(() => {
    if (!agent) return;
    if (!connected) return;
    const socket = getSocket();
    if (!socket) return;
    socket.on('agent:state:update', (newState: AgentState) => {
      if (agent?.status !== newState) {
        setStatus(newState);
        setTimeout(() => {
          setAgent(agent);

        }, 3000);
      }
    });
  }, [connected, agent]);

  return (
    <WrappedView>
      <BackButton />
      <YStack flex={1} items="center" justify="center">
        <YStack gap={16} mb={8} items="center">
          <AlarmSmoke size={64} color="$orange10" />
          <Text fontSize={32} fontWeight="bold">
            ¡Ups!
          </Text>
        </YStack>
        <YStack gap="$3" width="95%">
          <YStack>
            <Info>No tenemos información del agente.</Info>
            <Info>Por favor, para dar de alta deberás ejecutar lo siguiente:</Info>
          </YStack>
          <YStack gap="$1" bg="$gray2" p="$3" rounded="$2">
            <ScrollView horizontal showsHorizontalScrollIndicator={false}>
              <YStack gap="$2" shrink={0}>
                {cmds.map((cmd, i) => (
                  <XStack key={i} gap={4}>
                    <MonoText color="$gray10" opacity={0.5}>
                      -#
                    </MonoText>
                    <MonoText>{cmd}</MonoText>
                  </XStack>
                ))}
              </YStack>
            </ScrollView>
            <XStack
              position="absolute"
              b={4}
              r={4}
              pressStyle={{ opacity: 0.5 }}
              onPress={() => Clipboard.setStringAsync(cmds.join('\n'))}>
              <Copy size={16} opacity={0.5} />
            </XStack>
          </YStack>
          <Info fontSize={10} self="center">
            Después de ejecutarlo, el agente se registrará automáticamente.
          </Info>
        </YStack>
      </YStack>
      <YStack gap="$2" self="center" pb={32}>
        <Info color={status === 'ACTIVE' ? '$green10' : '$gray10'}>
          {status === 'ACTIVE' ? '¡Conexión detectada!' : 'Esperando una conexión...'}
        </Info>
        <Separator
          bg={status === 'ACTIVE' ? '$green10' : '$gray4'}
          p="$1"
          rounded="$2"
          width={180}
        />
      </YStack>
    </WrappedView>
  );
}
