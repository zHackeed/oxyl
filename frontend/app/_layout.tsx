import {
  useFonts,
  Inter_400Regular,
  Inter_600SemiBold,
  Inter_700Bold,
} from '@expo-google-fonts/inter';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { SplashScreen, Stack } from 'expo-router';
import { TamaguiProvider } from '@tamagui/core';
import tamaguiConfig from '@/tamagui.config';
import { useEffect, useState } from 'react';
import { DarkTheme, ThemeProvider } from '@react-navigation/native';
import { AuthStatus, useAuthStore } from '@/store/auth/useAuthStore';

SplashScreen.preventAutoHideAsync();

export default function RootLayout() {
  const [fontsLoaded] = useFonts({
    Inter: Inter_400Regular,
    InterSemiBold: Inter_600SemiBold,
    InterBold: Inter_700Bold,
  });

  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    if (fontsLoaded) {
      useAuthStore
        .getState()
        .hydrate()
        .then(() => {
          SplashScreen.hideAsync();
          setIsReady(true);
        });
    }
  }, [fontsLoaded]);

  if (!isReady) return null;

  const status = useAuthStore.getState().status;

  return (
    <SafeAreaProvider style={{ flex: 1 }}>
      <TamaguiProvider config={tamaguiConfig} defaultTheme="dark">
        <ThemeProvider value={DarkTheme}>
          <Stack screenOptions={{ headerShown: false }}>
            <Stack.Protected guard={status !== AuthStatus.AUTHENTICATED}>
              <Stack.Screen name="(auth)" />
            </Stack.Protected>
            <Stack.Protected guard={status === AuthStatus.AUTHENTICATED}>
              <Stack.Screen name="(company)" />
            </Stack.Protected>
          </Stack>
        </ThemeProvider>
      </TamaguiProvider>
    </SafeAreaProvider>
  );
}
