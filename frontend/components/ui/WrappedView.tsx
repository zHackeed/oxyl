import { Platform } from 'react-native';
import { KeyboardAvoidingView } from 'react-native-keyboard-controller';

import { SafeAreaView } from 'react-native-safe-area-context';
import { styled, View, ViewProps } from 'tamagui';

export const WrappedViewContainer = styled(View, {
  flex: 1,
  p: 10,
  pb: 20,
});

export const SafeAreaViewStyled = styled(SafeAreaView, {
  flex: 1,
  bg: '$background',
});

export function WrappedViewDismissable({ children, ...props }: ViewProps) {
  return (
    <SafeAreaViewStyled>
      <KeyboardAvoidingView
        style={{ flex: 1 }}
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}>
        <WrappedViewContainer {...props}>{children}</WrappedViewContainer>
      </KeyboardAvoidingView>
    </SafeAreaViewStyled>
  );
}

export function WrappedViewUnsafeDismissable({ children, ...props }: ViewProps) {
  return (
    <KeyboardAvoidingView
      style={{ flex: 1 }}
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}>
      <WrappedViewContainer {...props}>{children}</WrappedViewContainer>
    </KeyboardAvoidingView>
  );
}

export function WrappedView({ children, ...props }: ViewProps) {
  return (
    <SafeAreaViewStyled>
      <WrappedViewContainer {...props}>{children}</WrappedViewContainer>
    </SafeAreaViewStyled>
  );
}
