
import { XStack, Text, Separator } from "tamagui";

export interface InfoRowProps {
  label: string;
  value: string;
}

export function InfoRow({ label, value }: InfoRowProps) {
  return (
    <>
      <XStack paddingVertical="$3" justifyContent="space-between" alignItems="center">
        <Text color="$gray10" fontSize="$3">{label}</Text>
        <Text color="$white" fontSize="$3" fontWeight="500">{value}</Text>
      </XStack>
      <Separator borderColor="$gray6" />
    </>
  );
}