import React from 'react';
import { styled, View, YStack, Text } from 'tamagui';
import { CartesianChart, Line, Area, PointsArray } from 'victory-native';

const Container = styled(View, {
  flex: 1,
  bg: '$gray1',
  rounded: '$4',
  borderWidth: 0.5,
  borderColor: '$gray3',
overflow: 'hidden',
  shadowColor: '$black',
  shadowOffset: { width: 0, height: 2 },
  shadowOpacity: 0.1,
  shadowRadius: 4,
});

export const ChartContainer = styled(View, {
  mb: -8,
  ml: -8,
  mr: -8,
});

export const CHART_DEFAULTS = {
  xAxis: { lineWidth: 0 },
  padding: { left: 0, right: 0, top: 8, bottom: 0 },
  domainPadding: { left: 0, right: 0, top: 0, bottom: 0 },
} as const;

function BaseMetricCard({
  title,
  value,
  height = 100,
  legend,
  children,
  onPress,
}: {
  title: string;
  value: string;
  height?: number;
  children: React.ReactNode;
  legend?: React.ReactNode;
  onPress?: () => void;
}) {
  return (
    <Container onPress={onPress}>
      <YStack px={16} pt={16} pb={8}>
        <Text fontSize={12} color="$gray11" mb="$1">
          {title}
        </Text>
        <Text fontSize={24} fontWeight="600" color="#ffffff">
          {value}
        </Text>
      </YStack>
      <View height={height} mb={-8} ml={-8}>
        {children}
      </View>
      {legend && (
        <YStack px={16} pt={8} pb={16}>
          {legend}
        </YStack>
      )}
    </Container>
  );
}

export const MetricCard = React.memo(BaseMetricCard)
