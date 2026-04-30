import { YStack, Text } from "tamagui";

type EmptyStateProps = {
  icon: React.ReactNode;
  message: string;
  hint?: string;
};

export const EmptyState = ({ icon, message, hint }: EmptyStateProps) => (
  <YStack items="center" gap={16}>
    {icon}
    <YStack gap={2} items="center">
      <Text color="$red11">{message}</Text>
      {hint && (
        <Text color="$gray11" fontSize={12}>{hint}</Text>
      )}
    </YStack>
  </YStack>
);