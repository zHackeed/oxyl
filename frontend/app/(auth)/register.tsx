import { WrappedView } from '@/components/ui/WrappedView';
import { H2, YStack, Text, Button } from 'tamagui';
import { InputField } from '@/components/ui/InputField';    
import { useRouter } from 'expo-router';
import { useState } from 'react';
import { useAuthStore } from '@/store/auth/useAuthStore';
import { useAuthFacade } from '@/store/auth/useAuthFacade';
import { AuthService } from '@/lib/service/auth-service';

const Register = () => {
    const router = useRouter();

    const [name, setName] = useState("");
    const [surname, setSurname] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");


    const handleSignIn = async () => {
        if (name === '' || surname === '' || email === '' || password === '' || confirmPassword === '') {
            return;
        }

        // todo: validation logic
        // ->

        await AuthService.register(name, surname, email, password)
        router.back()
    }

    return (
        <WrappedView alignItems="center">
            <H2>Create a new account</H2>
            
            <YStack width="90%" gap={16} marginTop={32} >
                <YStack gap={8}>
                    <Text>Name</Text>
                    <InputField placeholder="Jhon" onChangeText={setName} />
                </YStack>
                 <YStack gap={8}>
                    <Text>Surname</Text>
                    <InputField placeholder="Doe"  onChangeText={setSurname}/>
                </YStack>
                 <YStack gap={8}>
                    <Text>Email</Text>
                    <InputField placeholder="Email" autoCapitalize="none" autoComplete="email" autoCorrect={false} onChangeText={setEmail} />
                </YStack>
                 <YStack gap={8}>
                    <Text>Password</Text>
                    <InputField placeholder="Enter password"  autoCapitalize="none" secureTextEntry size="$4" color="$white" onChangeText={setPassword} />
                </YStack>
                 <YStack gap={8}>
                    <Text>Confirm Password</Text>
                    <InputField placeholder="Confirm password"  autoCapitalize="none" secureTextEntry size="$4" color="$white" onChangeText={setConfirmPassword} />
                </YStack>
                
                <Button 
                    backgroundColor="$green9" 
                    width="75%" 
                    pressStyle={{ backgroundColor: '$green10' }} 
                    alignSelf="center" 
                    marginTop={16}
                    onPress={handleSignIn}
                >
                    Register
                </Button>
            </YStack>
        </WrappedView>
    );
}


export default Register;