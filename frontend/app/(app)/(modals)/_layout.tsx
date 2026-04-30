import { Stack } from 'expo-router';

export default function ModalLayout() {
  return (
    <Stack
      screenOptions={{
        headerShown: false,
        contentStyle: { backgroundColor: '#0e0e0e' },
      }}>
      <Stack.Screen name="new-company" />
      <Stack.Screen name="new-agent" />
    </Stack>
  );
}
