import { Stack } from "expo-router";

export default function AppLayout() {
  return (
    <Stack screenOptions={{ 
      headerShown: false, 
    }}>
      <Stack.Screen name="(selector)" />
      <Stack.Screen name="company/[id]" />
      <Stack.Screen name="agent/[id]" />
      <Stack.Screen
        name="(modals)"
        options={{
          presentation: 'formSheet',
          sheetCornerRadius: 24,
          sheetAllowedDetents: [0.75],
          sheetElevation: 30,
          sheetGrabberVisible: true,
          sheetShouldOverflowTopInset: false,
          contentStyle: { 
            height: '100%', 
          }, // there is a bug with formSheet when is nested. It freaks out and doesn't render properly.
        }}
      />
    </Stack>
  );
}