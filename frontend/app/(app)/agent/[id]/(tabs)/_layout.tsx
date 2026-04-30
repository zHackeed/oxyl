import { useAgentFarcade } from '@/store/agent/useAgentFarcade';
import { Icon, Label, NativeTabs } from 'expo-router/unstable-native-tabs';
import { useWebsocketFarcade } from '@/store/websocket/useWebsocketFarcade';
import { useEffect } from 'react';
import { getSocket } from '@/store/websocket/useWebsocketStore';
import { AgentState } from '@/lib/api/models/agent';

export default function AgentTabsLayout() {
  const { setStatus } = useAgentFarcade();
  const { connected } = useWebsocketFarcade();

  useEffect(() => {
    if (!connected) return;
    const websocket = getSocket();

    websocket?.on('agent:state:update', (newStatus: AgentState) => {
      setStatus(newStatus);
    });

    return () => {
      websocket?.off('agent:state:update');
    };
  }, [connected]);

  return (
    <NativeTabs labelStyle={{ color: '#e85d20' }} tintColor="#e85d20">
      <NativeTabs.Trigger name="index">
        <Label>Metricas</Label>
        <Icon sf="doc.text.magnifyingglass" />
      </NativeTabs.Trigger>
      <NativeTabs.Trigger name="notifications">
        <Label>Notificaciones</Label>
        <Icon sf="bell" />
      </NativeTabs.Trigger>
    </NativeTabs>
  );
}
