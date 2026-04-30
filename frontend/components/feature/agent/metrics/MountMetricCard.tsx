import { AgentMountPointMetric } from '@/lib/api/models/metrics';
import { useMemo } from 'react';
import { CHART_DEFAULTS, MetricCard } from '../BasetMetricCard';
import { CartesianChart, Line } from 'victory-native';
import { Text, View, XStack } from 'tamagui';

interface DiskMetricCardProps {
  data: Record<string, AgentMountPointMetric[]>;
  agentId: string;
}

const MOUNT_COLORS = ['#ef4444', '#FF7856', '#eab308', '#22c55e', '#a855f7'];
type DiskChartPoint = { timestamp: number } & Record<string, number>;

export function DiskChartCard({ data, agentId }: DiskMetricCardProps) {
  const mountKeys = useMemo(() => Object.keys(data), [data]);


  const chartData = useMemo<DiskChartPoint[]>(() => {
    const index = new Map<number, DiskChartPoint>();
    Object.entries(data).forEach(([mount, points]) => {
      points.forEach((p) => {
        const ts = new Date(p.when).getTime();
        if (!index.has(ts)) index.set(ts, { timestamp: ts });
        index.get(ts)![mount] = p.disk_usage / 1024 / 1024 / 1024;
      });
    });
    return Array.from(index.values()).sort((a, b) => a.timestamp - b.timestamp);
  }, [data]);

  const latest = chartData[chartData.length - 1];

  const legend = (
    <>
      {mountKeys.map((mount, i) => (
        <XStack key={mount} items="center" gap="$2">
          <View
            width={8}
            height={8}
            rounded="$10"
            bg={
              MOUNT_COLORS[
                i % MOUNT_COLORS.length
              ] as any /* I know this is a bad practice, but it's the only way to make it work and idc rn honeslty*/
            }
          />
          <Text fontSize={11} color="$gray10">
            {mount} {latest ? `${latest[mount]?.toFixed(1)} GB` : ''}
          </Text>
        </XStack>
      ))}
    </>
  );

  if (!mountKeys.length) return null;

  return (
    <MetricCard title="Espacio en uso" value={''} legend={legend}>
      <CartesianChart
        data={chartData}
        xKey="timestamp"
        yKeys={mountKeys}
        domain={{ x: [chartData[0]?.timestamp || Date.now(), Date.now()], y: [0, 50] }}
        yAxis={[{ lineWidth: 0 }]}
        {...CHART_DEFAULTS}>
        {({ points }) => (
          <>
            {mountKeys.map((mount, i) => (
              <Line
                key={mount}
                points={points[mount]}
                color={MOUNT_COLORS[i % MOUNT_COLORS.length]}
                strokeWidth={1.5}
                curveType="monotoneX"
                animate={{ type: 'timing', duration: 100 }}
              />
            ))}
          </>
        )}
      </CartesianChart>
    </MetricCard>
  );
}
