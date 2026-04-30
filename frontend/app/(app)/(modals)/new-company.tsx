import { BaseModal } from '@/components/feature/modals/BaseModal';
import { ModalEntry } from '@/components/feature/modals/ModalEntry';
import { CreateCompanyRequest } from '@/lib/api/requests/company';
import { companyService } from '@/lib/service/company';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useState } from 'react';
import { useRouter } from 'expo-router';
import { useCompanyFacade } from '@/store/company/userCompanyFacade';

export default function CreateCompanyModal() {
  const router = useRouter();
  const { setCompany } = useCompanyFacade();
  const [registationData, setRegistationData] = useState<CreateCompanyRequest>({
    display_name: '',
    webhook_type: 'DISCORD',
    webhook_url: '',
    channel: '',
  });

  const [errors, setErrors] = useState<Error | null>(null);
  const queryClient = useQueryClient();

  const registrationEntry = useMutation({
    mutationFn: () => companyService.createCompany(registationData),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['current-companies'] });
      setCompany(data.company);
      router.dismiss();
    },
    onError: (error) => {
      setErrors(error as Error);
    },
  });

  return (
    <BaseModal
      header="Nueva compañía"
      onSubmit={async () => await registrationEntry.mutate()}
      submitValue="Crear compañía"
      errors={errors}>
      <ModalEntry
        name="Nombre de la compañía"
        defaultValue="Compañia"
        consumeValue={(value) => {
          setRegistationData({ ...registationData, display_name: value });
        }}
      />
      <ModalEntry
        name="Tipo de webhook"
        defaultValue="Discord o Slack..."
        consumeValue={(value) => {
          setRegistationData({ ...registationData, webhook_type: value as 'DISCORD' | 'SLACK' });
        }}
      />
      <ModalEntry
        name="Punto de notificacion de alertas"
        defaultValue="https://discord.com/api/webhooks/..."
        consumeValue={(value) => {
          setRegistationData({ ...registationData, webhook_endpoint: value });
        }}
      />
      <ModalEntry
        name="Canal de notificaciones"
        defaultValue="123456789"
        consumeValue={(value) => {
          setRegistationData({ ...registationData, webhook_channel: value });
        }}
      />
    </BaseModal>
  );
}
