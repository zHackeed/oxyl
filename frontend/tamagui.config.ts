import { config } from '@tamagui/config/v3'
import { createFont, createTamagui } from '@tamagui/core'

const interFont = createFont({
  family: 'Inter',
  size: config.fonts.body.size,
  lineHeight: config.fonts.body.lineHeight,
  color: {
    1: 'white',
  },
  weight: {
    400: '400',
    600: '600',
    700: '700',
  },
  letterSpacing: config.fonts.body.letterSpacing,
  face: {
    400: { normal: 'Inter' },
    600: { normal: 'InterSemiBold' },
    700: { normal: 'InterBold' },
  },
})


const tamaguiConfig = createTamagui({
  ...config,
  fonts: {
    ...config.fonts,
    body: interFont,
  },
  themes: {
    ...config.themes,
    dark: {
      ...config.themes.dark,
      background: '#0e0e0e',
      backgroundHover: '#1a1a1a',
      backgroundPress: '#222222',
      backgroundStrong: '#111111',
      color: '#ffffff',
    }
  }
})

export type AppConfig = typeof tamaguiConfig
declare module '@tamagui/core' {
  interface TamaguiCustomConfig extends AppConfig {}
}

export default tamaguiConfig