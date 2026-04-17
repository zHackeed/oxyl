import { YStack, XStack, H2, Text, Separator } from 'tamagui';

// <Bell self="center" size={24} color="$orange8" mr="$2" />

export interface HeaderProps {
  title: string;
  description: string;
  icon?: React.ReactNode;
}

export default function GlobalHeader({ title, description, icon }: HeaderProps) {
  return (
    <YStack m="$3" gap={3}>
      <XStack gap="$2" justify="space-between">
        <H2>{title}</H2>
        {icon && icon}
      </XStack>

      <Text mt="$2" fontSize="$2" fontWeight={'400'} color="$color7" pb="$4">
        {description}
      </Text>
      <Separator borderColor="$gray12" mt="$2" mb="$2" />
    </YStack>
  );
}
