import { AgentNetworkMetric } from '@/lib/api/models/metrics';
import { useMemo } from 'react';
import { CHART_DEFAULTS, MetricCard } from '../BasetMetricCard';
import { Area, CartesianChart, Line } from 'victory-native';
import { XStack } from 'tamagui';

interface NetworkMetricCardProps {
  data: Record<string, AgentNetworkMetric[]>;
}

export function NetworkChartCard({ data }: NetworkMetricCardProps) {
  const ifaceKeys = useMemo(() => Object.keys(data), [data]);

  const rxData = useMemo(
    () =>
      (data[ifaceKeys[0]] ?? []).map((p) => ({
        timestamp: new Date(p.when).getTime(),
        value: p.rx_rate,
      })),
    [data, ifaceKeys]
  );

  const txData = useMemo(
    () =>
      (data[ifaceKeys[0]] ?? []).map((p) => ({
        timestamp: new Date(p.when).getTime(),
        value: p.tx_rate,
      })),
    [data, ifaceKeys]
  );

  const maxRx = Math.max(...rxData.map((d) => d.value), 1);
  const maxTx = Math.max(...txData.map((d) => d.value), 1);

  const latest = data[ifaceKeys[0]]?.[data[ifaceKeys[0]]?.length - 1];

  const formatRate = (bytes: number) =>
    bytes > 1024 * 1024
      ? `${(bytes / 1024 / 1024).toFixed(1)} MB/s`
      : bytes > 1024
        ? `${(bytes / 1024).toFixed(1)} KB/s`
        : `${bytes} B/s`;

  return (
    <XStack gap="$4">
      <MetricCard title="Red ↓" value={latest ? formatRate(latest.rx_rate) : '0 B/s'}>
        <CartesianChart
          data={rxData}
          xKey="timestamp"
          yKeys={['value']}
          domain={{
            y: [0, maxRx],
          }}
          yAxis={[{ lineWidth: 0 }]}
          {...CHART_DEFAULTS}>
          {({ points, chartBounds }) => (
            <>
              <Area
                points={points.value}
                y0={chartBounds.bottom}
                color="#7c22c5"
                opacity={0.18}
                curveType="monotoneX"
                animate={{ type: 'timing', duration: 100 }}
              />
              <Line
                points={points.value}
                color="#7c22c5"
                strokeWidth={1.5}
                curveType="monotoneX"
                animate={{ type: 'timing', duration: 100 }}
              />
            </>
          )}
        </CartesianChart>
      </MetricCard>
      <MetricCard title="Red ↑" value={latest ? formatRate(latest.tx_rate) : '0 B/s'}>
        <CartesianChart
          data={txData}
          xKey="timestamp"
          yKeys={['value']}
          domain={{
            y: [0, maxTx],
          }}
          yAxis={[{ lineWidth: 0 }]}
          {...CHART_DEFAULTS}>
          {({ points, chartBounds }) => (
            <>
              <Area
                points={points.value}
                y0={chartBounds.bottom}
                color="#7c22c5"
                opacity={0.18}
                curveType="monotoneX"
                animate={{ type: 'timing', duration: 100 }}
              />
              <Line
                points={points.value}
                color="#7c22c5"
                strokeWidth={1.5}
                curveType="monotoneX"
                animate={{ type: 'timing', duration: 100 }}
              />
            </>
          )}
        </CartesianChart>
      </MetricCard>
    </XStack>
  );
}
