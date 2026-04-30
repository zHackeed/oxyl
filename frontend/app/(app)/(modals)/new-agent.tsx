import { BaseModal } from '@/components/feature/modals/BaseModal';
import { ModalEntry } from '@/components/feature/modals/ModalEntry';
import { CreateAgentRequest } from '@/lib/api/requests/agent';
import { agentService } from '@/lib/service/agent';
import { useCompanyFacade } from '@/store/company/userCompanyFacade';
import { useMutation } from '@tanstack/react-query';
import { useRouter } from 'expo-router';
import { useState } from 'react';

export default function CreateNewAgent() {
  const { activeCompany } = useCompanyFacade();
  const router = useRouter();

  const [registerData, setRegisterData] = useState<CreateAgentRequest>({
    holder: activeCompany!.id,
    display_name: '',
    registered_ip: '',
  });
  const [errors, setErrors] = useState<Error | null>(null);

  const registerProcessor = useMutation({
    mutationFn: () => agentService.create(registerData),
    onSuccess: () => {
      router.dismiss();
    },
    onError: (error) => {
      setErrors(error);
    },
  });

  return (
    <BaseModal
      header="Nuevo agente"
      onSubmit={async () => {
        await registerProcessor.mutateAsync();
      }}
      submitValue="Crear agente"
      errors={errors}>
      <ModalEntry
        name="Nombre del agente"
        defaultValue=""
        consumeValue={(value) => {
          setRegisterData({
            ...registerData,
            display_name: value,
          });
        }}
      />
      <ModalEntry
        name="Ip del agente"
        defaultValue=""
        consumeValue={(value) => {
          setRegisterData({
            ...registerData,
            registered_ip: value,
          });
        }}
      />
    </BaseModal>
  );
}
