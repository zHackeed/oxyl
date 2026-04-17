import { InputField } from '@/components/ui/InputField';
import { YStack, styled, Text } from 'tamagui';

export interface ModalEntryProps {
  name: string;
  defaultValue: string;
  consumeValue: (value: string) => void;
}

const EntryContainer = styled(YStack, {
  gap: '$2',
});

const EntryLabel = styled(Text, {
  fontSize: 12,
  fontWeight: '600',
  letterSpacing: 1,
  color: '$color9',
  textTransform: 'uppercase',
});

export function ModalEntry({ name, defaultValue, consumeValue }: ModalEntryProps) {
  return (
    <EntryContainer>
      <EntryLabel>{name}</EntryLabel>
      <InputField placeholder={defaultValue} onChangeText={consumeValue} />
    </EntryContainer>
  );
}
