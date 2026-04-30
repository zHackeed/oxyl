import { WrappedView } from '@/components/ui/WrappedView';
import ModalRequest from '@/components/ui/ModalRequest';
import { Company } from '@/lib/api/models/company';
import { companyService } from '@/lib/service/company';
import { CompanyCard } from '@/components/feature/company/CompanyCard';
import { Link } from 'expo-router';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { FlatList } from 'react-native';
import GlobalHeader from '@/components/ui/Header';
import { useEffect } from 'react';
import { getSocket } from '@/store/websocket/useWebsocketStore';
import { useWebsocketFarcade } from '@/store/websocket/useWebsocketFarcade';

const CompaniesScreen = () => {
  const { connected: socketConnected } = useWebsocketFarcade();
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery<Company[], Error>({
    queryKey: ['current-companies'],
    queryFn: () => companyService.getCompanies(),
  });

  useEffect(() => {
    if (!socketConnected) return;
    const socket = getSocket();
    if (!socket) return;

    const handleAdd = (company: Company) => {
      queryClient.setQueryData<Company[]>(['current-companies'], (prev) => [...(prev ?? []), company]);
    };

    const handleRemove = (companyId: string) => {
      queryClient.setQueryData<Company[]>(['current-companies'], (prev) =>
        prev?.filter((c) => c.id !== companyId) ?? []
      );
    };

    socket.on('company:added', handleAdd);
    socket.on('company:removed', handleRemove);

    return () => {
      socket.off('company:added', handleAdd);
      socket.off('company:removed', handleRemove);
    };
  }, [queryClient, socketConnected]);

  const headerDescription = () => {
    if (isLoading) return 'Cargando...';
    if (error) return error.message;
    return 'Selecciona una compañía para continuar';
  };

  return (
    <WrappedView>
      <GlobalHeader title="Tus compañías" description={headerDescription()} />
      <FlatList
        data={data ?? []}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => <CompanyCard company={item} />}
        contentContainerStyle={{ gap: 12, paddingBottom: 140 }}
        showsVerticalScrollIndicator={false}
        scrollIndicatorInsets={{ bottom: 140 }}
        ListEmptyComponent={
          !isLoading ? (
            <Link href="/(app)/(modals)/new-company" push asChild>
              <ModalRequest
                title="¿Quieres registrar tu primera compañía?"
                subTitle="Haz clic aquí"
              />
            </Link>
          ) : null
        }
        ListFooterComponent={
          data && data.length > 0 ? (
            <Link href="/(app)/(modals)/new-company" push asChild>
              <ModalRequest title="¿Quieres registrar una compañía?" subTitle="Haz clic aquí" />
            </Link>
          ) : null
        }
      />
    </WrappedView>
  );
};

export default CompaniesScreen;