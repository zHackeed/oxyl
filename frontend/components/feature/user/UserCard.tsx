import React from "react";
import { View, Text, styled, YStack} from "tamagui";

interface UserCardProps {
  name: string;
  surname: string;
  creation_date: string;
}

const Container = styled(YStack, {
  gap: 16,
  items: "center",
  rounded: "$4",
  pt: "$4",
})

const ViewCard = styled(View, {
  bg: "$gray4",
  width: "$10",
  height: "$10",
  rounded: "$4",
  items: "center",
  justify: "center",
  borderWidth: 1,
  borderColor: "$black7",
})

export function UserCard({ name, surname, creation_date }: UserCardProps) {
  const initials = `${name[0]}${surname[0]}`;
  return (
    <Container>
      <ViewCard>
        <Text fontSize={32} color="$white">{initials}</Text> 
      </ViewCard>
      <Text fontSize={16} color="$gray8">Miembro desde el {new Date(creation_date).toLocaleDateString()}</Text>
    </Container>
    
  );
}