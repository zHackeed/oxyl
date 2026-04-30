import { WrappedView } from '@/components/ui/WrappedView';
import { useAuthFacade } from '@/store/auth/useAuthFacade';
import { Button, Text, YStack, Spinner, H2, Separator } from 'tamagui';
import { User } from '@/lib/api/models/user';
import { userService } from '@/lib/service/user';
import { UserCard } from '@/components/feature/user/UserCard';
import { InputField } from '@/components/ui/InputField';
import { useQuery } from '@tanstack/react-query';

export default function Account() {
  const { signOut } = useAuthFacade();
  const { data, isLoading } = useQuery<User | null>({
    queryKey: ['user'],
    queryFn: () => userService.get(),
    refetchOnWindowFocus: true,
  });

  return (
    <WrappedView px="$4" pt="$4">
      <YStack pb="$4">
        <H2 mt="$4" fontWeight="800">
          Tu perfil
        </H2>
        <Text mt="$2" fontSize="$2" fontWeight={'400'} color="$color7">
          Esta es la información de tu cuenta.
        </Text>
      </YStack>

      <YStack gap="$6">
        {isLoading ? (
          <YStack flex={1} items="center" justify="center">
            <Spinner size="large" color="$green9" />
          </YStack>
        ) : data ? (
          <>
            <UserCard name={data.name} surname={data.surname} creation_date={data.created_at} />
            <Separator width="75%" borderColor="$gray12" self="center" />
            <YStack gap="$4">
              <YStack gap="$2">
                <Text
                  fontSize={12}
                  fontWeight="600"
                  letterSpacing={1}
                  color="$color9"
                  textTransform="uppercase">
                  Nombre
                </Text>
                <InputField value={data.name} disabled />
              </YStack>
              <YStack gap="$2">
                <Text
                  fontSize={12}
                  fontWeight="600"
                  letterSpacing={1}
                  color="$color9"
                  textTransform="uppercase">
                  Apellido
                </Text>
                <InputField value={data.surname} disabled />
              </YStack>
              <YStack gap="$2">
                <Text
                  fontSize={12}
                  fontWeight="600"
                  letterSpacing={1}
                  color="$color9"
                  textTransform="uppercase">
                  Correo electrónico
                </Text>
                <InputField placeholder={data.email} disabled />
              </YStack>
            </YStack>
          </>
        ) : (
          <YStack flex={1} items="center" justify="center">
            <Text color="$color9">No se pudo cargar la información de tu cuenta.</Text>
          </YStack>
        )}

        <YStack gap="$3" pb="$6">
          <Button size="$4" bg="$orange9" borderColor="$orange9" borderWidth={1} onPress={signOut}>
            Cerrar sesión
          </Button>
        </YStack>
      </YStack>
    </WrappedView>
  );
}
