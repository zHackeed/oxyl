import GlobalHeader from '@/components/ui/Header';
import { WrappedView } from '@/components/ui/WrappedView';
import { useAgentFarcade } from '@/store/agent/useAgentFarcade';
import { useQuery } from '@tanstack/react-query';
import { NotificationCard } from '@/components/feature/agent/notification/NotificationCard';
import { Text } from 'tamagui';
import { agentService } from '@/lib/service/agent';
import { AgentNotificationLog } from '@/lib/api/models/agent';
import { FlatList } from 'react-native';
import { EmptyState } from '@/components/ui/EmptyState';
import { Ghost } from '@tamagui/lucide-icons-2';

export default function Notifications() {
  const { agent } = useAgentFarcade();
  const { data, isLoading } = useQuery<AgentNotificationLog[]>({
    queryKey: ['agent-notifications', agent!.id],
    queryFn: () => agentService.fetchNotifications(agent!.id),
  });

  return (
    <WrappedView>
      <GlobalHeader title="Notificaciones" description={`${data?.length || 0} notificaciones`} />
      <FlatList
        data={data ?? []}
        renderItem={({ item }) => <NotificationCard notification={item} />}
        keyExtractor={(item) => item.identifier}
        contentContainerStyle={{ gap: 12 }}
        showsVerticalScrollIndicator={false}
        ListEmptyComponent={
          !isLoading ? (
            <EmptyState
              icon={<Ghost size={32} color="$red6" />}
              message="No hay notificaciones para este agente"
              hint="Si crees que esto es un error, contacta con tu administrador"
            />
          ) : (
            <Text>Cargando...</Text>
          )
        }
      />
    </WrappedView>
  );
}
