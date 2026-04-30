import { WrappedView } from '@/components/ui/WrappedView';
import { Text, ScrollView, YStack } from 'tamagui';
import ThresholdCard from '@/components/feature/company/ThresholdCard';
import { useQuery } from '@tanstack/react-query';
import { ActiveCompanyThreshold, ThresholdMetadataMap } from '@/lib/api/models/company';
import { companyService } from '@/lib/service/company';
import GlobalHeader from '@/components/ui/Header';
import { useCompanyFacade } from '@/store/company/userCompanyFacade';

const Thresholds = () => {
  const { activeCompany } = useCompanyFacade();

  const { data, isLoading, isLoadingError } = useQuery<ActiveCompanyThreshold[]>({
    queryKey: ['active-company-thresholds', activeCompany?.id],
    queryFn: () => companyService.getActiveCompanyThresholds(activeCompany?.id || ''),
    refetchOnWindowFocus: true,
  });

  const headerDescription = () => {
    if (isLoading) return 'Cargando...';
    if (isLoadingError) return 'Error al cargar los umbrales.';
    return 'Configura los umbrales que activarán notificaciones automáticas';
  };

  const validThresholds = data
    ?.filter((t) => ThresholdMetadataMap[t.threshold_id])
    .sort((a, b) => a.threshold_id.localeCompare(b.threshold_id));

  return (
    <WrappedView>
      <GlobalHeader
        title="Tus Umbrales"
        description={headerDescription()}
      />
      {!isLoading && !isLoadingError && (
        <ScrollView rounded="$7" mb="$8">
          <YStack gap="$4">
            {validThresholds?.map((threshold) => (
              <ThresholdCard
                key={threshold.threshold_id}
                type={threshold.threshold_id}
                value={threshold.value}
                limit={ThresholdMetadataMap[threshold.threshold_id].limit}
              />
            ))}
          </YStack>
        </ScrollView>
      )}
    </WrappedView>
  );
};

export default Thresholds;