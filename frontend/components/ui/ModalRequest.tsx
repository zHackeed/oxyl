import { View, Text, styled } from 'tamagui';
import { CirclePlus } from '@tamagui/lucide-icons-2';
import { ViewProps } from 'react-native-svg/lib/typescript/fabric/utils';
import { GestureResponderEvent } from 'react-native';

export interface ModalRequestProps extends ViewProps {
  title: string;
  subTitle: string;
  onPress?: (event: GestureResponderEvent) => void;
}

const ModalView = styled(View, {
  mt: 'auto',
  height: '$14',
  width: '100%',
  rounded: '$3',
  borderWidth: 2,
  borderColor: '$color9',
  borderStyle: 'dashed',
  items: 'center',
  justify: 'center',
  self: 'center',
  pressStyle: {
    opacity: 0.7,
    scale: 0.98,
  },
});

export default function ModalRequest({ title, subTitle, onPress }: ModalRequestProps) {
  return (
    <ModalView onPress={onPress}>
      <CirclePlus size={32} marginBottom="$3" color="$color9" />
      <Text color="$color8" mb="$3" fontWeight={400}>
        {title}
      </Text>
      <Text color="$orange9" fontWeight={300}>
        {subTitle}
      </Text>
    </ModalView>
  );
}
