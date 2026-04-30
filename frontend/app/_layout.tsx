import '@tamagui/native/setup-burnt';
import tamaguiConfig from '@/components/ui/tamagui.config';
import { SplashScreen, Stack } from 'expo-router';
import { useAuthFacade } from '@/store/auth/useAuthFacade';
import { AuthStatus } from '@/store/auth/useAuthStore';
import { SafeAreaProvider, initialWindowMetrics } from 'react-native-safe-area-context';
import { KeyboardProvider } from 'react-native-keyboard-controller';
import { TamaguiProvider, Theme } from '@tamagui/core';
import * as SystemUI from 'expo-system-ui';
import React, { useEffect } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useWebsocketStore } from '@/store/websocket/useWebsocketStore';
import { AppState } from 'react-native';

SplashScreen.preventAutoHideAsync();
SystemUI.setBackgroundColorAsync('black'); // So for whatever reason, it is ignoring the background field on my app.json. This fixes it and makes it bearable.

const queryClient = new QueryClient();

export default function RootLayout() {
  return (
    <SafeAreaProvider initialMetrics={initialWindowMetrics} style={{ flex: 1 }}>
      <RootLayourInternal />
    </SafeAreaProvider>
  );
}

function RootLayourInternal() {
  const { status } = useAuthFacade();
  const { connect, disconnect } = useWebsocketStore();

  useEffect(() => {
    if (status !== AuthStatus.LOADING) {
      console.log('Hiding splash screen');
      SplashScreen.hideAsync();
    }
  }, [status]);

  useEffect(() => {
    if (status !== AuthStatus.AUTHENTICATED) {
      disconnect();
      return;
    }

    connect();

    const subscription = AppState.addEventListener('change', (nextState) => {
      if (nextState === 'active') connect();
      else if (nextState === 'background' || nextState === 'inactive') disconnect();
    });

    return () => subscription.remove();
  }, [status]);

  if (status === AuthStatus.LOADING) {
    console.log('Loading...');
    return null;
  }

  return (
    <TamaguiProvider config={tamaguiConfig} defaultTheme="dark">
      <QueryClientProvider client={queryClient}>
        <Theme name="dark">
          <KeyboardProvider>
            <Stack
              screenOptions={{
                headerShown: false,
                contentStyle: { backgroundColor: '#101010' },
              }}>
              <Stack.Protected guard={status === AuthStatus.UNAUTHENTICATED}>
                <Stack.Screen name="(auth)" />
              </Stack.Protected>
              <Stack.Protected guard={status === AuthStatus.AUTHENTICATED}>
                <Stack.Screen name="(app)" />
              </Stack.Protected>
            </Stack>
          </KeyboardProvider>
        </Theme>
      </QueryClientProvider>
    </TamaguiProvider>
  );
}
