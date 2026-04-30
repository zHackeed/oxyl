import { AgentGeneralMetric } from "@/lib/api/models/metrics";
import { useMemo } from "react";
import { CHART_DEFAULTS, MetricCard } from "../BasetMetricCard";
import { Area, CartesianChart, Line } from "victory-native";

interface CpuMetricCardProps {
  data: AgentGeneralMetric[];
}

export function CpuChartCard({ data }: CpuMetricCardProps) {
  const chartData = useMemo(() =>
    data.map((item) => ({
      timestamp: new Date(item.when).getTime(),
      value: item.cpu_usage,
    })), [data]);

  const latest = data[data.length - 1];

  return (
    <MetricCard title="CPU" value={`${latest?.cpu_usage?.toFixed(2) || 0}%`}>
      <CartesianChart
        data={chartData}
        xKey="timestamp"
        yKeys={['value']}
        domain={{ x: [chartData[0]?.timestamp || Date.now(), Date.now()], y: [0, 100] }}
        yAxis={[{ lineWidth: 0 }]}
        {...CHART_DEFAULTS}
      >
        {({ points, chartBounds }) => (
          <>
            <Area
              points={points.value}
              y0={chartBounds.bottom}
              color="#00ffc8"
              opacity={0.18}
              curveType="monotoneX"
              animate={{ type: 'timing', duration: 100 }}
            />
            <Line
              points={points.value}
              color="#00ffc8"
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