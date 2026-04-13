import { Button, Form, styled } from 'tamagui'

export const PrimaryButton = styled(Button, {
  size: "$5",
  width: "75%", // 75% of the view width
  bg: "$orange9",
  self: "center",
  pressStyle: {
    opacity: 0.7,
    scale: 0.98,
  }
})

export function SubmitterButton(props: any) {
  return (
    <Form.Trigger asChild>
      <PrimaryButton {...props} />
    </Form.Trigger>
  )
}