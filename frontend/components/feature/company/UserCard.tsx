import { Badge } from '@/components/ui/Badge';
import { CompanyMember, isAdmin } from '@/lib/api/models/company';
import { View, XStack, styled, Text, YStack } from 'tamagui';

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

const Initials = styled(View, {
  bg: '$color4',
  rounded: '$3',
  items: 'center',
  justify: 'center',
  width: 48,
  height: 48,
});

export interface CompanyUserCardProps {
  member: CompanyMember;
  onPress?: () => void;
}

export function CompanyUserCard({ member, onPress }: CompanyUserCardProps) {
  const initials = `${member.user.name.at(0)?.toUpperCase() || 'U'}${member.user.surname.at(0)?.toUpperCase() || 'U'}`;
  return (
    <StyledXStack onPress={onPress}>
      <Initials>
        <Text fontSize="$5">{initials}</Text>
      </Initials>
      <YStack gap={7} flex={1}>
        <XStack gap={8}>
          <Text fontWeight="bold" fontSize={18} flex={1}>
            {member.user.name} {member.user.surname}
          </Text>

          {isAdmin(member) ? (
            <Badge borderColor="$orange8" bg="$orange6">
              Admin
            </Badge>
          ) : (
            <Badge borderColor="$blue8" bg="$blue6">
              Miembro
            </Badge>
          )}
        </XStack>

        <Text color="$gray11">{member.user.email}</Text>
      </YStack>
    </StyledXStack>
  );
}
