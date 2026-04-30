import GlobalHeader from '@/components/ui/Header';
import { WrappedView } from '@/components/ui/WrappedView';
import { useQuery } from '@tanstack/react-query';
import { companyService } from '@/lib/service/company';
import { FlatList } from 'react-native';
import { CompanyUserCard } from '@/components/feature/company/UserCard';
import { useCompanyFacade } from '@/store/company/userCompanyFacade';

// Todo: Add add and remove functionality

const Users = () => {
  const { activeCompany } = useCompanyFacade();

  const { data, isLoading, isLoadingError } = useQuery({
    queryKey: ['active-company-users', activeCompany?.id],
    queryFn: () => companyService.getMembers(activeCompany?.id || ''),
  });

  const headerDescription = () => {
    if (isLoading) return 'Cargando...';
    if (isLoadingError) return 'Error al cargar los miembros.';
    return 'Gestiona los miembros de tu compañía';
  };

  return (
    <WrappedView>
      <GlobalHeader title="Miembros" description={headerDescription()} />
      {!isLoading && !isLoadingError && (
        <FlatList
          data={data}
          renderItem={({ item }) => <CompanyUserCard member={item} />}
          keyExtractor={(item) => item.user.email}
          contentContainerStyle={{ gap: 12 }}
          showsVerticalScrollIndicator={false}
        />
      )}
    </WrappedView>
  );
};

export default Users;