import { WrappedView } from '@/components/ui/WrappedView';
import { H2, Text, YStack, Separator } from 'tamagui';
import ModalRequest from '@/components/ui/ModalRequest';
import { Company } from '@/lib/api/models/company';
import { companyService } from '@/lib/service/company';
import { CompanyCard } from '@/components/feature/company/CompanyCard';
import { Link } from 'expo-router';
import { useQuery } from '@tanstack/react-query';
import { FlatList } from 'react-native';
import GlobalHeader from '@/components/ui/Header';

const CompaniesScreen = () => {
  const { data, isLoading, error } = useQuery<Company[], Error>({
    queryKey: ['current-companies'],
    queryFn: () => companyService.getCompanies(),
    refetchOnWindowFocus: true,
  });

  return (
    <WrappedView>
      <YStack>
        <GlobalHeader title="Tus compañías" description="Selecciona una compañía para continuar" />
        {isLoading ? (
          <Text mt="$2" fontSize="$2" fontWeight={'400'} color="$color7">
            Cargando...
          </Text>
        ) : data && data.length > 0 ? (
          <>
            {data.length > 0 && (
              <FlatList
                data={data}
                renderItem={({ item }) => <CompanyCard key={item.id} company={item} />}
                style={{ borderRadius: 8 }}
                contentContainerStyle={{ gap: 12, paddingBottom: 140 }}
                showsVerticalScrollIndicator={false}
                keyExtractor={(item) => item.id}
                scrollIndicatorInsets={{ bottom: 140 }}
                ListFooterComponent={
                  <Link href="/(modals)/new-company" push asChild>
                    <ModalRequest
                      title="¿Quieres registrar una compañia?"
                      subTitle="Haz clic aquí"
                    />
                  </Link>
                }
              />
            )}

            {!data && (
              <Link href="/(modals)/new-company" push asChild>
                <ModalRequest
                  title="¿Quieres registrar tu primera compañia?"
                  subTitle="Haz clic aquí"
                />
              </Link>
            )}
          </>
        ) : error ? (
          <Text mt="$2" fontSize="$2" fontWeight={'400'} color="$color7">
            {error.message}
          </Text>
        ) : null}
      </YStack>
    </WrappedView>
  );
};

export default CompaniesScreen;
