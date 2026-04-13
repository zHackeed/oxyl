import { ChevronLeft } from '@tamagui/lucide-icons-2';
import { router } from 'expo-router';
import { XStack, Text } from 'tamagui';

export interface BackButtonProps {
  onPress?: () => void;
}

export function BackButton({ onPress }: BackButtonProps) {
  return (
    <XStack
      onPress={() => {
        console.log("presed")
        if (onPress) {
          onPress();
        }
        router.back();
      }}
      gap="$1"
      items="center">
      <ChevronLeft color="#FF7856" size={24} />
      <Text color="#FF7856">Volver</Text>
    </XStack>
  );
}
