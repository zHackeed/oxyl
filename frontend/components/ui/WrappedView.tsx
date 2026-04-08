import { styled, View, ViewProps} from 'tamagui';
import { Pressable, Keyboard } from 'react-native';

export const WrappedViewContainer = styled(View, {
    flex: 1,
    padding: 20,
    justifyContent: 'center',
    backgroundColor: '$background',
})

export function WrappedViewDismissable({children, ...props }: ViewProps) {
  return (
    <Pressable style={{ flex: 1 }} onPress={Keyboard.dismiss}>
        <WrappedViewContainer {...props}>
            {children}
        </WrappedViewContainer>
    </Pressable>
  )
}

export function WrappedView({children, ...props }: ViewProps) {
  return (
    <WrappedViewContainer {...props}>
      {children}
    </WrappedViewContainer>
  )
}
