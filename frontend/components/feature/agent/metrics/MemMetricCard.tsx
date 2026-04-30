import { AgentGeneralMetric } from '@/lib/api/models/metrics';
import { useMemo } from 'react';
import { CHART_DEFAULTS, MetricCard } from '../BasetMetricCard';
import { Area, CartesianChart, Line } from 'victory-native';
import { chartDomain } from '@/lib/utils/charts';

interface MemMetricCardProps {
  data: AgentGeneralMetric[];
  totalMemory: number;
}

export function MemChartCard({ data, totalMemory }: MemMetricCardProps) {
  const chartData = useMemo(
    () =>
      data.map((item) => ({
        timestamp: new Date(item.when).getTime(),
        value: item.memory_usage / 1024,
      })),
    [data]
  );

  const latest = data[data.length - 1];
  const usedMb = latest ? (latest.memory_usage / 1024 / 1024).toFixed(1) : '0';
  const totalMb = (totalMemory / 1024).toFixed(0);

  return (
    <MetricCard title="Memoria" value={`${usedMb} MB / ${totalMb} MB`}>
      <CartesianChart
        data={chartData}
        xKey="timestamp"
        yKeys={['value']}
        domain={{
          y: [0, totalMemory],
        }}
        yAxis={[{ lineWidth: 0 }]}
        {...CHART_DEFAULTS}>
        {({ points, chartBounds }) => (
          <>
            <Area
              points={points.value}
              y0={chartBounds.bottom}
              color="#eeff00"
              opacity={0.18}
              curveType="monotoneX"
              animate={{ type: 'timing', duration: 100 }}
            />
            <Line
              points={points.value}
              color="#eeff00"
              strokeWidth={1.5}
              curveType="monotoneX"
              animate={{ type: 'timing', duration: 100 }}
            />
          </>
        )}
      </CartesianChart>
    </MetricCard>
  );
}
