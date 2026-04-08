import { Link } from "expo-router";
import { InputField } from "@/components/ui/InputField";
import { WrappedViewDismissable } from "@/components/ui/WrappedView";
import { useEffect, useState } from "react";
import { YStack, Text, Button } from "tamagui";
// @ts-ignore
import LogoSvg from "@/assets/logo-notext.svg";
import { useAuthFacade } from "@/store/auth/useAuthFacade";

export default function SignIn() {
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');

  const { signIn } = useAuthFacade()

  useEffect(() => {
    if (email || password) {
      setEmail('');
      setPassword('');
    }
  }, [email, password]);

  const handleSignIn = () => {
    if (email === '' || password === '') {
      return
    }

    signIn(email, password)
  };


  return (
    <WrappedViewDismissable>
      <YStack
        width={300}
        alignSelf="center"
        alignItems="center"
      >
        <LogoSvg width={140} height={140} />
        <Text fontSize="$8" fontWeight="bold">Welcome to Oxyl</Text>

        <YStack gap={20} width="100%" margin={20}>
            <YStack gap={6} width="100%">
              <Text color="white" fontSize="$4">Email</Text>
              <InputField placeholder="john.doe@example.com" onChangeText={setEmail} autoCapitalize="none" autoComplete='email' />
            </YStack>

            <YStack gap={6} width="100%">
              <Text color="white" fontSize="$4">Password</Text>
              <InputField secureTextEntry placeholder="••••••••••••••••••••" onChangeText={setPassword} autoCapitalize="none" autoComplete='password' autoCorrect={false} />
            </YStack>
            
            { /*<Link href="/forgot" style={{ alignSelf: 'flex-end', marginBottom: 10 }}>
              <Text color="gray" fontSize="$2">Forgot your password?</Text>
            </Link> */}
        </YStack>          
    
        <YStack gap={20} width="100%" alignItems="center" marginTop={20}>
          <Button backgroundColor="$orange9" width="75%" pressStyle={{
            backgroundColor: '$orange12',
          }} onPress={handleSignIn}>
            Sign In
          </Button>

          <Text color="gray" fontSize="$4" alignSelf='center'>Don't have an account? <Link href="/register"><Text color="$orange9">Register</Text></Link></Text>
        </YStack>
      </YStack>
    </WrappedViewDismissable>
  );
}
