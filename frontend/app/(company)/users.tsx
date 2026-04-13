import GlobalHeader from '@/components/ui/Header';
import { WrappedView } from '@/components/ui/WrappedView';
import { useQuery } from '@tanstack/react-query';
import { companyService } from '@/lib/service/company';
import { useCompanyFacade } from '@/store/company/userCompanyFacade';
import { ScrollView, Text } from 'tamagui';
import { FlatList } from 'react-native';
import { CompanyUserCard } from '@/components/feature/company/UserCard';

// Todo: Add add and remove functionality

const Users = () => {
  const { activeCompany } = useCompanyFacade();

  const { data, isLoading, isLoadingError } = useQuery({
    queryKey: ['active-company-users'],
    queryFn: () => companyService.getMembers(activeCompany?.id || ''),
  });

  return (
    <WrappedView>
      <GlobalHeader title="Miembros" description="Gestiona los miembros de tu compañía" />
      
      {isLoading ? (
        <Text>Cargando...</Text>
      ) : isLoadingError ? (
        <Text>Error al cargar los miembros</Text>
      ) : (
        <FlatList
          data={data}
          renderItem={({ item }) => {
            return <CompanyUserCard member={item} />;
          }}
          keyExtractor={(item) => item.user.email}
        />
      )}
    </WrappedView>
  );
};

export default Users;
