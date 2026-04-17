import { Company } from '@/lib/api/models/company';
import { useCompanyFacade } from '@/store/company/userCompanyFacade';
import { Building2, ChevronRight } from '@tamagui/lucide-icons-2';
import { useRouter } from 'expo-router';
import { View, Text, XStack, YStack, styled } from 'tamagui';

export interface CompanyCardProps {
  company: Company;
}

const StyledXStack = styled(XStack, {
  bg: '$color2',
  p: '$4',
  rounded: '$4',
  items: 'center',
  gap: '$3',
  borderWidth: 1,
  borderColor: '#2a2a2a',
  pressStyle: {
    scale: 0.98,
    bg: '$color3',
  },
});

const IconStyle = styled(View, {
  p: '$3',
  bg: '$color4',
  rounded: '$3',
  borderColor: '#4a4a4a',
  borderWidth: 1,
});

export function CompanyCard({ company }: CompanyCardProps) {
  const router = useRouter();
  const { setCompany } = useCompanyFacade();
  return (
    <StyledXStack
      onPress={() => {
        setCompany(company);
        router.push('/(company)');
      }}>
      <IconStyle>
        <Building2 size={16} color="$color8" />
      </IconStyle>

      <YStack flex={1} gap="$1">
        <Text fontSize="$5" fontWeight="400" color="$color12">
          {company.display_name}
        </Text>
        <XStack items="center" gap="$1">
          <View bg="$green9" width={6} height={6} rounded="$10" mr="$1" />
          <Text fontSize="$2" color="$color9">
            X / {company.limit_nodes} agents
          </Text>
        </XStack>
      </YStack>

      <ChevronRight size={18} color="$color8" />
    </StyledXStack>
  );
}
