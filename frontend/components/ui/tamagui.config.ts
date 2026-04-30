import { defaultConfig } from '@tamagui/config/v5';
import { createTamagui } from '@tamagui/core';

const tamaguiConfig = createTamagui({
  ...defaultConfig,
  themes: {
    ...defaultConfig.themes,
    dark: {
      ...defaultConfig.themes.dark,
      background: '#0e0e0e',
      backgroundHover: '#1a1a1a',
      backgroundPress: '#222222',
      backgroundStrong: '#111111',
      color: '#ffffff',
      secondary: '#a0a0a0',
    },
  },
});

export type AppConfig = typeof tamaguiConfig;
declare module '@tamagui/core' {
  interface TamaguiCustomConfig extends AppConfig {}
}

export default tamaguiConfig;