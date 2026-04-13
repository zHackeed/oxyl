import { WrappedViewUnsafeDismissable } from '@/components/ui/WrappedView';
import { YStack, Text, Separator, Form, Spinner } from 'tamagui';
import { useEffect, useState } from 'react';
import { SubmitterButton } from '@/components/ui/Button';
import { useRouter } from 'expo-router';

export interface BaseModalProps {
  header: string;
  submitValue: string;
  children: React.ReactNode;
  onSubmit: () => Promise<void>;
  errors: Error | null;
}

export function BaseModal({ header, submitValue, children, onSubmit, errors }: BaseModalProps) {
  const [submitting, setSubmitting] = useState(false);
  const [localErrors, setLocalErrors] = useState<Error | null>(null);


  useEffect(() => {
    if (submitting == true) {
      const submitting = async () => {
        try {
          await onSubmit();
        } finally {
          setSubmitting(false);
        }
      };
      submitting();
    }
  }, [submitting]);
  
  useEffect(() => {
    setLocalErrors(errors);
  }, [errors]);

  return (
    <WrappedViewUnsafeDismissable px="$5" pt="$6" pb="$8">
      <Form onSubmit={() => setSubmitting(true)} flex={1}>
        <YStack gap="$1" mb="$6" mt="$2" justify="center" items="center">
          <Text fontSize={32} fontWeight="700" color="$color12">
            {header}
          </Text>
          <Separator borderColor="#2a2a2a" m="$2" width="$10" />
          <Text fontSize={14} self="center" color="$color9">
            Rellena los siguientes campos
          </Text>
        </YStack>
        <Separator borderColor="#2a2a2a" mb="$6" />

        <YStack gap="$5">{children}</YStack>

        {submitting && <Spinner color="$yellow12" mt="auto" self="center" size="large" />}
        {localErrors && (
          <Text color="$red12" mt="$2">
            {localErrors.message}
          </Text>
        )}
        <Form.Trigger asChild>
          <SubmitterButton mt="auto" self="center" disabled={submitting}>
            {submitValue}
          </SubmitterButton>
        </Form.Trigger>
      </Form>
    </WrappedViewUnsafeDismissable>
  );
}
