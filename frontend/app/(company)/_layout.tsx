import { CompanyPermission, hasPermission } from '@/lib/api/models/company';
import { companyService } from '@/lib/service/company';
import { useCompanyFacade } from '@/store/company/userCompanyFacade';
import { useQuery } from '@tanstack/react-query';
import { useRouter } from 'expo-router';
import { Icon, Label, NativeTabs } from 'expo-router/unstable-native-tabs';
import React, { useEffect } from 'react';

export default function CompanyLayout() {
  const router = useRouter();
  const { permissions, setCompanyPermissions, activeCompany } = useCompanyFacade();

  const { data, isLoading, isLoadingError } = useQuery({
    queryKey: ['company', activeCompany?.id],
    queryFn: () => companyService.getPermission(activeCompany!.id),
  });

  useEffect(() => {
    if (isLoadingError) {
      router.back();
      return;
    }

    if (data) {
      setCompanyPermissions(data.permissions);
    }
  }, [data, isLoadingError, isLoading, router, activeCompany, setCompanyPermissions]);

  if (isLoading) return null;

  return (
    <NativeTabs labelStyle={{ color: '#e85d20' }} tintColor="#e85d20">
      <NativeTabs.Trigger name="index">
        <Label>Agentes</Label>
        <Icon sf="desktopcomputer" drawable="custom_android_drawable" />
      </NativeTabs.Trigger>
      <NativeTabs.Trigger
        name="users"
        hidden={!hasPermission(permissions ?? [], CompanyPermission.ManageMembers)}>
        <Label>Usuarios</Label>
        <Icon sf="person.fill" drawable="custom_android_drawable" />
      </NativeTabs.Trigger>
      <NativeTabs.Trigger
        name="thresholds"
        hidden={!hasPermission(permissions ?? [], CompanyPermission.ManageThresholds)}>
        <Label>Umbrales</Label>
        <Icon sf="chart.bar.fill" drawable="custom_android_drawable" />
      </NativeTabs.Trigger>
    </NativeTabs>
  );
}
