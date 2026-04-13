import { Agent } from '@/lib/api/models/agent';
import { useCompanyFacade } from '@/store/company/userCompanyFacade';
import { useQuery } from '@tanstack/react-query';
import { agentService } from '@/lib/service/agent';
import { FlatList } from 'react-native';
import { WrappedView } from '@/components/ui/WrappedView';
import { BackButton } from '@/components/ui/BackButton';
import GlobalHeader from '@/components/ui/Header';
import { Text, YStack } from 'tamagui';
import { Link } from 'expo-router';
import ModalRequest from '@/components/ui/ModalRequest';
import AgentCard from '@/components/feature/agent/AgentCard';

const AgentIndex = () => {
  const { setCompany, activeCompany, setCompanyPermissions } = useCompanyFacade();
  const { data, isLoading } = useQuery<Agent[] | null>({
    queryKey: ['company-agents'],
    queryFn: () => {
      return agentService.get(activeCompany?.id || '');
    },
  });

  const hasAgents = data && data.length > 0;

  return (
    <WrappedView>
      <BackButton
        onPress={() => {
          setCompany(null);
          setCompanyPermissions(null);
        }}
      />
      <GlobalHeader
        title="Tus Agentes"
        description={
          isLoading
            ? 'Cargando...'
            : data?.length === 0
              ? 'No hay agentes actualmente registrados en la compañía'
              : `En la compañía ${activeCompany?.display_name} hay un total de ${data?.length} agentes`
        }
      />

      <FlatList
        data={data ?? []}
        renderItem={({ item }) => {
          return <AgentCard name={item.display_name} description={item.registered_ip} />;
        }}
        keyExtractor={(item) => item.id}
        contentContainerStyle={{ gap: 12 }}
        showsVerticalScrollIndicator={false}
        ListEmptyComponent={
          <Link href="/(modals)/new-agent" push asChild>
            <ModalRequest title="¿Quieres registrar tu primer agente?" subTitle="Haz clic aquí" />
          </Link>
        }
        ListFooterComponent={
          !hasAgents ? null : (
            <Link href="/(modals)/new-agent" push asChild>
              <ModalRequest title="¿Quieres registrar un agente?" subTitle="Haz clic aquí" />
            </Link>
          )
        }
      />
    </WrappedView>
  );
};

export default AgentIndex;
