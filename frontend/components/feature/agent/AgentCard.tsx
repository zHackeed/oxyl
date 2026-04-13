import { View, XStack, YStack, styled, Text } from 'tamagui';
import { Server, ChevronRight } from '@tamagui/lucide-icons-2';

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
  name: string;
  description: string;
}

export default function AgentCard({ name, description }: AgentCardProps) {
  return (
    <Container>
      <XStack gap={12} items="center">
        <IconContainer>
          <Server size={24} color="$color" />
        </IconContainer>
        <YStack gap={4} justify="center" flex={1}>
          <Text fontWeight="bold" fontSize={18}>
            {name}
          </Text>
          <Text color="$gray11">{description}</Text>
        </YStack>
        <ChevronRight size={18} color="$color8" />
      </XStack>

      {/* todo: websocket impl */}
      <View
        height={100}
        bg="$color4"
        flex={1}
        width="100%"
        rounded="$4"
        mt="$2"
        justify="center"
        items="center">
        <Text color="$color10">Todo: websocket</Text>
      </View>
    </Container>
  );
}
