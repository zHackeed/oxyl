import { Icon, Label, NativeTabs } from 'expo-router/unstable-native-tabs';
import React from 'react';

export default function CompanyLayout() {
  return (
    <NativeTabs labelStyle={{ color: '#e85d20' }} tintColor="#e85d20">
      <NativeTabs.Trigger name="index">
        <Label>Compañías</Label>
        <Icon sf="house.fill" drawable="custom_android_drawable" />
      </NativeTabs.Trigger>
      <NativeTabs.Trigger name="account">
        <Label>Tu cuenta</Label>
        <Icon sf="person.fill" drawable="custom_android_drawable" />
      </NativeTabs.Trigger>
      <NativeTabs.Trigger name="new_company" hidden={true}>
        <Label>New company</Label>
        <Icon sf="plus.circle" drawable="custom_android_drawable" />
      </NativeTabs.Trigger>
    </NativeTabs>
  );
}
