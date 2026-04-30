import GlobalHeader from '@/components/ui/Header';
import { WrappedView } from '@/components/ui/WrappedView';
import { useAgentFarcade } from '@/store/agent/useAgentFarcade';
import { YStack } from 'tamagui';
import { InfoRow } from '@/components/feature/agent/InfoRow';
import { XCircle } from '@tamagui/lucide-icons-2';
import { useRouter } from 'expo-router';

export default function AgentInfo() {
  const { dismiss } = useRouter();
  const { agent } = useAgentFarcade();

  if (!agent) return null;

  return (
    <WrappedView>
      <XCircle
        size={24}
        onPress={() => dismiss()}
        position="absolute"
        t={24}
        l={24}
        color="$orange9"
        border="$4"
        p="$3"
      />
      <GlobalHeader
        title="Información del agente"
        description={agent.display_name + ' - ' + agent.registered_ip}
      />

      <YStack gap="$3">
        <InfoRow label="ID" value={agent.id} />
        <InfoRow label="Sistema operativo" value={agent.metadata!.system_os ?? '—'} />
        <InfoRow label="Modelo CPU" value={agent.metadata!.cpu_model ?? '—'} />
        <InfoRow
          label="Memoria total"
          value={agent.metadata!.total_memory ? formatBytes(agent.metadata!.total_memory) : '—'}
        />
        <InfoRow
          label="Disco total"
          value={agent.metadata!.total_disk ? formatBytes(agent.metadata!.total_disk) : '—'}
        />
        <InfoRow
          label="Último handshake"
          value={agent.last_handshake ? formatDate(agent!.last_handshake) : '—'}
        />
      </YStack>
    </WrappedView>
  );
}

// Nice
function formatBytes(bytes: number, decimals = 2) {
  if (!+bytes) return '0 Bytes';

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
}

function formatDate(date: number): string {
  return new Date(date).toLocaleString('es-ES', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}
