import { CompanyPermission, hasPermission } from '@/lib/api/models/company';
import { companyService } from '@/lib/service/company';
import { useCompanyFacade } from '@/store/company/userCompanyFacade';
import { useWebsocketFarcade } from '@/store/websocket/useWebsocketFarcade';
import { useQuery } from '@tanstack/react-query';
import { useRouter, useLocalSearchParams } from 'expo-router';
import { Icon, Label, NativeTabs } from 'expo-router/unstable-native-tabs';
import React, { useEffect } from 'react';

export default function CompanyLayout() {
  const router = useRouter();
  const { id } = useLocalSearchParams<{ id: string }>();
  const { setCompany, setPermissions, permissions, activeCompany } = useCompanyFacade();
  const { connected: websocketConnected, join, leave } = useWebsocketFarcade();

  const { data: company } = useQuery({
    queryKey: ['company', id],
    queryFn: () => companyService.getCompany(id || activeCompany?.id || ''),
    initialData: activeCompany
  });

  const { data: permsData, isLoadingError } = useQuery({
    queryKey: ['company-permissions', id],
    queryFn: () => companyService.getPermission(id || activeCompany?.id || ''),
    initialData: { permissions: [] } as { permissions: CompanyPermission[] },
  });

  useEffect(() => {
    if (company) setCompany(company);
  }, [company]);

  useEffect(() => {
    if (!permsData) return;
    setPermissions(permsData.permissions);
  }, [permsData]);

  useEffect(() => {
    if (!websocketConnected) return;
    join('company', id);

    return () => {
      leave('company', id);
    };
  }, [websocketConnected]);

  useEffect(() => {
    if (isLoadingError) router.back();
  }, [isLoadingError]);

  useEffect(() => {
    return () => {
      setCompany(null);
      setPermissions(null);
    };
  }, []);

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
