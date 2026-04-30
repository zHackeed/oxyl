import GlobalHeader from '@/components/ui/Header';
import { WrappedView } from '@/components/ui/WrappedView';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { companyService } from '@/lib/service/company';
import { FlatList } from 'react-native';
import { CompanyUserCard } from '@/components/feature/company/UserCard';
import { useCompanyFacade } from '@/store/company/userCompanyFacade';
import { Alert } from 'react-native';
import { userService } from '@/lib/service/user';

const Users = () => {
  const { activeCompany } = useCompanyFacade();
  const queryClient = useQueryClient();

  const { data: currentUser } = useQuery({
    queryKey: ['user'],
    queryFn: () => userService.get(),
  });

  const { data, isLoading, isLoadingError } = useQuery({
    queryKey: ['active-company-users', activeCompany?.id],
    queryFn: () => companyService.getMembers(activeCompany?.id || ''),
  });

  const headerDescription = () => {
    if (isLoading) return 'Cargando...';
    if (isLoadingError) return 'Error al cargar los miembros.';
    return 'Gestiona los miembros de tu compañía';
  };

  const removeMemberMutation = useMutation({
    mutationFn: (email: string) => companyService.removeMember(activeCompany?.id || '', email),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ['active-company-users', activeCompany?.id],
      });
    },
    onError: (_) => {
      Alert.alert('Error', 'No se pudo eliminar el miembro');
    },
  });

  return (
    <WrappedView>
      <GlobalHeader title="Miembros" description={headerDescription()} />
      {!isLoading && !isLoadingError && (
        <FlatList
          data={data}
          renderItem={({ item }) => (
            <CompanyUserCard
              member={item}
              onPress={() => {
                if (item.user.email === currentUser?.email) {
                  Alert.alert('No puedes eliminarte a ti mismo');
                  return;
                }
                Alert.alert(
                  '¿Quieres eliminar este miembro?',
                  'Esta acción no se puede deshacer.',
                  [
                    { text: 'Cancelar' },
                    {
                      text: 'Eliminar',
                      onPress: () => removeMemberMutation.mutate(item.user.email),
                    },
                  ]
                );
              }}
            />
          )}
          keyExtractor={(item) => item.user.email}
          contentContainerStyle={{ gap: 12 }}
          showsVerticalScrollIndicator={false}
        />
      )}
    </WrappedView>
  );
};

export default Users;
