import { styled } from '@tamagui/core';
import { Input } from 'tamagui';

export const InputField = styled(Input, {
  backgroundColor: '#222222',
  borderWidth: 1,
  borderColor: '#333333',
  width: '100%',
  placeholder: 'Example text',
  shadowColor: '#000000',
  shadowOffset: { width: 0, height: 2 },
  shadowOpacity: 0.25,
  shadowRadius: 3.84,
});