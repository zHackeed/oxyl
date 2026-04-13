import { useRouter } from 'expo-router';
import { useState } from 'react';
import { UserRegisterRequest } from '@/lib/api/requests/user';
import { AuthService } from '@/lib/service/auth';
import { WrappedViewDismissable } from '@/components/ui/WrappedView';
import { H2, YStack, Text, Form } from 'tamagui';
import { InputField } from '@/components/ui/InputField';
import { SubmitterButton } from '@/components/ui/Button';

// TODO: Use modal instead of full view?
const Register = () => {
  const router = useRouter();

  const [message, setMessage] = useState<string>('');
  const [errors, setErrors] = useState<string>('');
  const [registationData, setRegistationData] = useState<UserRegisterRequest>({
    name: '',
    surname: '',
    email: '',
    password: '',
    confirmPassword: '',
  });

  const setFormData = (field: keyof UserRegisterRequest, value: string) => {
    setRegistationData({
      ...registationData,
      [field]: value,
    });
    setErrors('');
  };

  const handleSignIn = async () => {
    await AuthService.register(registationData)
      .then(async (success: boolean) => {
        if (success) {
          setMessage('Registrado, por favor inicia sesión...');
          setTimeout(() => router.back(), 3000);
        }
      })
      .catch((error) => {
        setErrors(error.message);
      });
  };

  return (
    <WrappedViewDismissable>
      <Form flex={1} onSubmit={handleSignIn} justify="center" items="center">
        <YStack>
          <H2>Crea una cuenta nueva</H2>
          <YStack gap={16} mt={32}>
            <YStack gap={8}>
              <Text fontFamily="$body">Nombre</Text>
              <InputField placeholder="Jhon" onChangeText={(value) => setFormData('name', value)} />
            </YStack>
            <YStack gap={8}>
              <Text fontFamily="$body">Apellido</Text>
              <InputField
                placeholder="Doe"
                onChangeText={(value) => setFormData('surname', value)}
              />
            </YStack>
            <YStack gap={8}>
              <Text fontFamily="$body">Email</Text>
              <InputField
                placeholder="Email"
                autoCapitalize="none"
                autoComplete="email"
                autoCorrect={false}
                onChangeText={(value) => setFormData('email', value)}
              />
            </YStack>
            <YStack gap={8}>
              <Text>Contraseña</Text>
              <InputField
                placeholder="Enter password"
                autoCapitalize="none"
                secureTextEntry
                size="$4"
                color="$color"
                onChangeText={(value) => setFormData('password', value)}
              />
            </YStack>
            <YStack gap={8}>
              <Text>Confirmar Contraseña</Text>
              <InputField
                placeholder="Confirmar contraseña"
                autoCapitalize="none"
                secureTextEntry
                size="$4"
                color="$color"
                onChangeText={(value) => setFormData('confirmPassword', value)}
              />
            </YStack>

            <YStack gap={8}>
              {errors && <Text color="red">{errors}</Text>}
              {message && <Text color="green">{message}</Text>}
            </YStack>
            
            <SubmitterButton bg="$green9">Registrarse</SubmitterButton>
          </YStack>
        </YStack>
      </Form>
    </WrappedViewDismissable>
  );
};

export default Register;
