import { CompanyThresholdNotificationType, ThresholdMetadataMap } from '@/lib/api/models/company';
import { companyService } from '@/lib/service/company';
import { useCompanyFacade } from '@/store/company/userCompanyFacade';
import { useMutation } from '@tanstack/react-query';
import { useRef, useState } from 'react';
import { H4, Text, XStack, YStack, Slider, styled, GetThemeValueForKey, Spinner } from 'tamagui';

const Container = styled(YStack, {
  rounded: '$3',
  borderWidth: 1,
  borderColor: '$gray11',
  bg: '$color2',
  width: '95%',
  self: 'center',
  p: '$4',
  gap: '$3',
});

export interface ThresholdCardProps {
  type: CompanyThresholdNotificationType;
  value: number;
  limit?: number
}

export default function ThresholdCard({ type, value, limit }: ThresholdCardProps) {
  const { activeCompany } = useCompanyFacade();
  const saveTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const [currentValue, setCurrentValue] = useState<number>(value);
  const [updating, setUpdating] = useState(false);

  const thresholdMetadata = ThresholdMetadataMap[type];

  const updateValue = useMutation({
    mutationFn: () => companyService.updateThreshold(activeCompany!.id, type, currentValue),
    onSuccess: () => {
      setUpdating(false);
    },
    onError: (error) => {
      console.error('Error updating threshold:', error);
      setUpdating(false);
    },
  });

  if (!thresholdMetadata) {
    return null;
  }

  const handleValueChange = ([v]: number[]) => {
    setCurrentValue(v);
    if (saveTimer.current) {
      clearTimeout(saveTimer.current);
    }
    setUpdating(true);
    saveTimer.current = setTimeout(() => {
      updateValue.mutate();
    }, 500);
  };

  return (
    <Container>
      <XStack justify="space-between" items="flex-start">
        <YStack gap="$1" pb="$2">
          <H4>{thresholdMetadata.label}</H4>
          <Text fontSize="$2" color="$color7">
            {thresholdMetadata.description}
          </Text>
        </YStack>
        {updating && <Spinner size="small" />}
      </XStack>

      <Text color="$color10" fontWeight="500">
        {currentValue} {limit ? `/${limit}` : '%'}
      </Text>
      <Slider defaultValue={[value]} onValueChange={handleValueChange} max={limit || 100} step={1}>
        <Slider.Track backgroundColor="$gray9">
          <Slider.TrackActive backgroundColor={thresholdMetadata.color} />
        </Slider.Track>
        <Slider.Thumb
          size={18}
          circular
          borderColor="$white"
          bg="$white"
          pressStyle={{
            opacity: 0.7,
            scale: 0.98,
            bg: '$gray11',
          }}
        />
      </Slider>
    </Container>
  );
}
