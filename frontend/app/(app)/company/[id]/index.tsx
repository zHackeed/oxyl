import { Agent, AgentState } from '@/lib/api/models/agent';
import { useCompanyFacade } from '@/store/company/userCompanyFacade';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { agentService } from '@/lib/service/agent';
import { FlatList } from 'react-native';
import { WrappedView } from '@/components/ui/WrappedView';
import GlobalHeader from '@/components/ui/Header';
import { Link  } from 'expo-router';
import ModalRequest from '@/components/ui/ModalRequest';
import AgentCard from '@/components/feature/agent/AgentCard';
import { getSocket } from '@/store/websocket/useWebsocketStore';
import { useEffect } from 'react';
import { CompanyPermission, hasPermission } from '@/lib/api/models/company';
import { Ghost } from '@tamagui/lucide-icons-2';
import { EmptyState } from '@/components/ui/EmptyState';
import { useWebsocketFarcade } from '@/store/websocket/useWebsocketFarcade';
import { BackButton } from '@/components/ui/BackButton';

const AgentIndex = () => {
  const { connected: socketConnected } = useWebsocketFarcade();
  const { activeCompany, permissions } = useCompanyFacade();
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery<Agent[] | null>({
    queryKey: ['company-agents', activeCompany?.id],
    queryFn: () => agentService.get(activeCompany?.id || ''),
  });

  useEffect(() => {
    if (!socketConnected) return;
    const socket = getSocket();
    if (!socket) return;

    const handleCreation = (agent: Agent) => {
      console.log(agent)
      queryClient.setQueryData<Agent[]>(['company-agents', activeCompany?.id], (prev) => [
        ...(prev ?? []),
        agent,
      ]);
    };

    const handleDeletion = (agentId: string) => {
      queryClient.setQueryData<Agent[]>(['company-agents', activeCompany?.id], (prev) =>
        prev?.filter((agent) => agent.id !== agentId)
      );
    };

    const handleStateUpdate = (agentId: string, newState: AgentState) => {
      queryClient.setQueryData<Agent[]>(
        ['company-agents', activeCompany?.id],
        (prev) => prev?.map((a) => (a.id === agentId ? { ...a, status: newState } : a)) ?? []
      );
    };

    socket.on('company:agent:creation', handleCreation);
    socket.on('company:agent:deletion', handleDeletion);
    socket.on('company:agent:update', handleStateUpdate);

    return () => {
      socket.off('company:agent:creation', handleCreation);
      socket.off('company:agent:deletion', handleDeletion);
      socket.off('company:agent:update', handleStateUpdate);
    };
  }, [socketConnected, queryClient, activeCompany?.id]);

  const headerDescription = () => {
    if (isLoading) return 'Cargando...';
    if (!data?.length) return 'No hay agentes actualmente registrados en la compañía';
    return `En la compañía ${activeCompany?.display_name} hay un total de ${data.length} agentes`;
  };

  if (!permissions) return null;

  return (
    <WrappedView>
      <BackButton />
      <GlobalHeader title="Tus Agentes" description={headerDescription()} />
      <FlatList
        data={data ?? []}
        renderItem={({ item }) => (
          <AgentCard agent={item} />
        )}
        keyExtractor={(item) => item.id}
        contentContainerStyle={{ gap: 12 }}
        showsVerticalScrollIndicator={false}
        ListEmptyComponent={
          !isLoading ? (
            hasPermission(permissions, CompanyPermission.ManageAgents) ? (
              <Link href="/(app)/(modals)/new-agent" push asChild>
                <ModalRequest title="¿Quieres registrar tu primer agente?" subTitle="Haz clic aquí" />
              </Link>
            ) : (
              <EmptyState
                icon={<Ghost size={32} color="$red6" />}
                message="No hay agentes registrados para esta compañía"
                hint="Si crees que esto es un error, contacta con tu administrador"
              />
            )
          ) : null
        }
        ListFooterComponent={
          hasPermission(permissions, CompanyPermission.ManageAgents) && data && data.length > 0 ? (
            <Link href="/(app)/(modals)/new-agent" push asChild>
              <ModalRequest title="¿Quieres registrar un agente?" subTitle="Haz clic aquí" />
            </Link>
          ) : null
        }
      />
    </WrappedView>
  );
};

export default AgentIndex;