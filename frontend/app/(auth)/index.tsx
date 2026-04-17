import { Link } from 'expo-router';
import { InputField } from '@/components/ui/InputField';
import { WrappedViewDismissable } from '@/components/ui/WrappedView';
import { useState } from 'react';
import { YStack, Text } from 'tamagui';
// @ts-ignore
import LogoSvg from '@/assets/logo-notext.svg';
import { useAuthFacade } from '@/store/auth/useAuthFacade';
import { UserLoginRequest } from '@/lib/api/requests/user';
import { usLoginFormSchema } from '@/lib/validators/authentication';
import { SubmitterButton } from '@/components/ui/Button';

export default function SignIn() {
  const [errors, setErrors] = useState('');
  const [loginData, setLoginData] = useState<UserLoginRequest>({
    email: '',
    password: '',
  });
  const { signIn } = useAuthFacade();

  const setFormData = (field: keyof UserLoginRequest, value: string) => {
    setLoginData({
      ...loginData,
      [field]: value,
    });
    setErrors('');
  };

  const handleSignIn = async () => {
    usLoginFormSchema
      .validate(loginData)
      .then(async () => {
        const success = await signIn(loginData.email, loginData.password);
        if (!success) {
          setErrors('Invalid email or password, please check your credentials');
        }
      })
      .catch((error) => {
        setErrors(error.message);
      });
  };

  return (
    <WrappedViewDismissable justify="center" items="center">
      <YStack width={'75%'} self="center" items="center">
        <LogoSvg width={140} height={140} />
        <Text fontSize="$8" fontWeight="bold" fontFamily="$body">
          Bienvenido a Oxyl
        </Text>

        <YStack gap={20} width="100%" m={10}>
          <YStack gap={6}>
            <Text fontSize="$4" fontFamily="$body">
              Email
            </Text>
            <InputField
              placeholder="john.doe@example.com"
              onChangeText={(text) => {
                setFormData('email', text);
              }}
              autoCapitalize="none"
              autoComplete="email"
            />
          </YStack>

          <YStack gap={6} width="100%">
            <Text fontSize="$4" fontFamily="$body">
              Password
            </Text>
            <InputField
              secureTextEntry
              placeholder="••••••••••••••••••••"
              onChangeText={(text) => {
                setFormData('password', text);
              }}
              autoCapitalize="none"
              autoComplete="password"
              autoCorrect={false}
            />
          </YStack>
        </YStack>

        {errors && (
          <Text color="red" fontSize="$2">
            {errors}
          </Text>
        )}

        <YStack gap={20} width="100%" items="center" mt={20}>
          <SubmitterButton onPress={handleSignIn}>Iniciar sesión</SubmitterButton>

          <Text color="gray" fontSize="$4" self="center" fontFamily="$body">
            No tienes cuenta?{' '}
            <Link href="/register">
              <Text color="$orange9">Registrate</Text>
            </Link>
          </Text>
        </YStack>
      </YStack>
    </WrappedViewDismissable>
  );
}
