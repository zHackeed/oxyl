import { XStack, YStack, Text } from 'tamagui';
import { AgentNotificationLog, notificationTypeLabels } from '@/lib/api/models/agent';
import React from 'react';

export type NotificationCardProps = {
  notification: AgentNotificationLog;
};

export function NotificationCard({ notification }: NotificationCardProps) {
  return (
    <YStack bg="$color2" p="$4" rounded="$4" borderWidth={1} borderColor="$gray4" gap="$2">
      <XStack gap={12} items="center">
        <YStack flex={1} gap={2}>
          <Text fontWeight="bold" fontSize={15}>
            {notificationTypeLabels[notification.trigger_reason]}
          </Text>
          <Text color="$gray11" fontSize={13}>
            Valor: {notification.trigger_value}
          </Text>
        </YStack>
        <Text color="$gray10" fontSize={12}>
          {new Date(notification.sent_at).toLocaleString('es-ES', {
            day: '2-digit', month: '2-digit',
            hour: '2-digit', minute: '2-digit',
          })}
        </Text>
      </XStack>
    </YStack>
  );
}

export default React.memo(NotificationCard);