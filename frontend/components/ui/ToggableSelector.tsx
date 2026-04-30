import { ToggleGroup, Text } from 'tamagui';

export interface ToggleSelectorOption {
  value: string;
  label: string;
}

interface ToggleSelectorProps {
  value: string;
  options: ToggleSelectorOption[];
  onValueChange: (value: string) => void;
}

export function ToggleSelector({ value, options, onValueChange }: ToggleSelectorProps) {
  return (
    <ToggleGroup
      mt="$1"
      mr="$2"
      type="single"
      value={value}
      onValueChange={onValueChange}
      disableDeactivation
      self="center"
      orientation="horizontal"
      p={4}
      gap="$1"
      flexDirection="row">
      {options.map(({ value: key, label }) => (
        <ToggleGroup.Item
          key={key}
          value={key}
          rounded="$2"
          px="$4"
          ml="$2"
          py="$2"
          minWidth="22%"
          justifyContent="center"
          pressStyle={{ bg: '$gray5' }}>
          <Text fontSize={12} color={value === key ? '$white' : '$white8'}>
            {label}
          </Text>
        </ToggleGroup.Item>
      ))}
    </ToggleGroup>
  );
}