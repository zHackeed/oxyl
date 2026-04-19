import tamaguiConfig from '@/components/ui/tamagui.config';
import { SplashScreen, Stack } from 'expo-router';
import { useAuthFacade } from '@/store/auth/useAuthFacade';
import { AuthStatus } from '@/store/auth/useAuthStore';
import {
  SafeAreaProvider,
  initialWindowMetrics,
  useSafeAreaInsets,
} from 'react-native-safe-area-context';
import { KeyboardProvider } from 'react-native-keyboard-controller';
import { TamaguiProvider, Theme } from '@tamagui/core';
import * as SystemUI from 'expo-system-ui';
import { useEffect, useState } from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useWebsocketStore } from '@/store/websocket/useWebsocketStore';
import { AppState } from 'react-native';

SplashScreen.preventAutoHideAsync();
SystemUI.setBackgroundColorAsync('black'); // So for whatever reason, it is ignoring the background field on my app.json. This fixes it and makes it bearable.

const queryClient = new QueryClient();

export default function RootLayout() {
  const { status } = useAuthFacade();
  const insets = useSafeAreaInsets();
  const { connect, disconnect } = useWebsocketStore();

  useEffect(() => {
    const subscription = AppState.addEventListener('change', (nextState) => {
      if (status !== AuthStatus.AUTHENTICATED) return;

      if (nextState === 'active') {
        connect();
      } else if (nextState === 'background' || nextState === 'inactive') {
        disconnect();
      }
    });

    return () => subscription.remove();
  }, [status]);

  useEffect(() => {
    if (status !== AuthStatus.LOADING) {
      SplashScreen.hideAsync();
    }

    if (status === AuthStatus.AUTHENTICATED) {
      connect();
    } else {
      disconnect();
    }

    return () => {
      if (status === AuthStatus.AUTHENTICATED) {
        disconnect();
      }
    };
  }, [status]);

  if (status === AuthStatus.LOADING) {
    return null;
  }

  return (
    <TamaguiProvider config={tamaguiConfig} defaultTheme="dark" insets={insets}>
      <SafeAreaProvider initialMetrics={initialWindowMetrics} style={{ flex: 1 }}>
        <QueryClientProvider client={queryClient}>
          <Theme name="dark">
            <KeyboardProvider>
              <Stack
                screenOptions={{
                  headerShown: false,
                  contentStyle: { backgroundColor: 'transparent' },
                }}>
                <Stack.Protected guard={status === AuthStatus.UNAUTHENTICATED}>
                  <Stack.Screen name="(auth)" />
                </Stack.Protected>
                <Stack.Protected guard={status === AuthStatus.AUTHENTICATED}>
                  <Stack.Screen name="(selector)" />
                  <Stack.Screen name="(company)" />
                </Stack.Protected>

                {/* modals */}
                <Stack.Screen
                  name="(modals)"
                  options={{
                    presentation: 'formSheet',
                    sheetCornerRadius: 24,
                    sheetAllowedDetents: [0.8],
                    sheetGrabberVisible: false,
                    sheetShouldOverflowTopInset: true,
                    contentStyle: {
                      height: '100%', // there is a bug with formSheet when is nested. It freaks out and doesn't render properly.
                    },
                  }}
                />
              </Stack>
            </KeyboardProvider>
          </Theme>
        </QueryClientProvider>
      </SafeAreaProvider>
    </TamaguiProvider>
  );
}
