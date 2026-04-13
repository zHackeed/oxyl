import { WrappedView } from '@/components/ui/WrappedView';
import { H2, Text, Separator, ScrollView, YStack, XStack } from 'tamagui';
import ThresholdCard from '@/components/feature/company/ThresholdCard';
import { useQuery } from '@tanstack/react-query';
import { ActiveCompanyThreshold, ThresholdMetadataMap } from '@/lib/api/models/company';
import { companyService } from '@/lib/service/company';
import { useCompanyFacade } from '@/store/company/userCompanyFacade';
import { Bell } from '@tamagui/lucide-icons-2';
import GlobalHeader from '@/components/ui/Header';

const Thresholds = () => {
  const { activeCompany } = useCompanyFacade();

  const { data, isLoading, isLoadingError } = useQuery<ActiveCompanyThreshold[]>({
    queryKey: ['active-company-thresholds'],
    queryFn: () => companyService.getActiveCompanyThresholds(activeCompany!.id),
    refetchOnWindowFocus: true,
  });

  return (
    <WrappedView>
      <GlobalHeader
        title="Tus Umbrales"
        description="Configura los umbrales que activarán notificaciones automáticas"
        icon={<Bell size={24} color="$orange8" />}
      />

      {isLoading ? (
        <Text fontSize="$2" color="$color7">
          Cargando...
        </Text>
      ) : isLoadingError ? (
        <Text fontSize="$2" color="$red10">
          Error al cargar los umbrales.
        </Text>
      ) : (
        <ScrollView rounded="$7">
          <YStack gap="$4">
            {data
              ?.sort((a, b) => a.threshold_id.localeCompare(b.threshold_id))
              .map((threshold) => {
                if (!ThresholdMetadataMap[threshold.threshold_id]) {
                  return null;
                }
                return (
                  <ThresholdCard
                    key={threshold.threshold_id}
                    type={threshold.threshold_id}
                    value={threshold.value}
                  />
                );
              })}
          </YStack>
        </ScrollView>
      )}
    </WrappedView>
  );
};

export default Thresholds;
