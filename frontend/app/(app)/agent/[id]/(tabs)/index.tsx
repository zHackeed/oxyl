import GlobalHeader from '@/components/ui/Header';
import { useAgentFarcade } from '@/store/agent/useAgentFarcade';
import { WrappedView } from '@/components/ui/WrappedView';
import { useState } from 'react';
import { BackButton } from '@/components/ui/BackButton';
import { ToggleSelector } from '@/components/ui/ToggableSelector';
import { useAgentMetrics } from '@/hooks/useAgentMetrics';
import { CpuChartCard } from '@/components/feature/agent/metrics/CpuMetricCard';
import { MemChartCard } from '@/components/feature/agent/metrics/MemMetricCard';
import { DiskChartCard } from '@/components/feature/agent/metrics/MountMetricCard';
import { NetworkChartCard } from '@/components/feature/agent/metrics/NetworkMetricCard';
import { ScrollView, YStack } from 'tamagui';
import { Info } from '@tamagui/lucide-icons-2';
import { useRouter } from 'expo-router';

const intervalLabels = {
  '15m': '15 mins',
  '1h': '1 hora',
  '6h': '6 horas',
  '7d': '7 dias',
} as const;

type IntervalKey = keyof typeof intervalLabels;

const options = Object.entries(intervalLabels).map(([key, label]) => ({
  value: key,
  label,
}));

export default function Agent() {
  const { push } = useRouter();
  const { agent, status } = useAgentFarcade();
  const [interval, setInterval] = useState<IntervalKey>('15m');
  const { general, network, mounts } = useAgentMetrics(agent?.id || '', interval);

  if (!agent) {
    return null;
  }

  return (
    <WrappedView>
      <BackButton />
      <GlobalHeader
        title={agent?.display_name || ''}
        description={
          status === 'ACTIVE'
            ? `Métricas actualizadas ${new Date(general.at(-1)?.when || 0).toLocaleString()}`
            : `Última vez visto: ${new Date(agent?.last_handshake || 0).toLocaleString()}`
        }
        icon={
          <Info
            size={24}
            color="$orange8"
            onPress={() => {
              console.log('Navigating to info page');
              push({
                pathname: '/agent/[id]/info',
                params: {
                  id: agent!.id,
                },
              });
            }}
          />
        }
      />
      <ScrollView mb={64} showsVerticalScrollIndicator={false}>
        <ToggleSelector
          options={options}
          value={interval}
          onValueChange={(value) => setInterval(value as IntervalKey)}
        />
        <YStack gap="$4" mt="$4">
          <CpuChartCard data={general} />
          <MemChartCard data={general} totalMemory={agent?.metadata?.total_memory || 0} />
          <DiskChartCard data={mounts} agentId={agent?.id || ''} />
          <NetworkChartCard data={network} />
        </YStack>
      </ScrollView>
    </WrappedView>
  );
}
