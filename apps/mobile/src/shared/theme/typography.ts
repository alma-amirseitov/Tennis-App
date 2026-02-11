export const fontSize = {
  xs: 11,
  sm: 12,
  base: 13,
  md: 14,
  lg: 15,
  xl: 16,
  '2xl': 18,
  '3xl': 20,
  '4xl': 22,
  '5xl': 24,
  '6xl': 28,
  '7xl': 32,
  '8xl': 36,
} as const;

export const fontWeight = {
  regular: '400' as const,
  medium: '500' as const,
  semibold: '600' as const,
  bold: '700' as const,
  extrabold: '800' as const,
};

// React Native expects lineHeight in absolute pixels, not as a multiplier
export const lineHeight = {
  tight: undefined, // let RN calculate automatically
  normal: undefined,
  relaxed: undefined,
} as const;

export const textStyles = {
  h1: { fontSize: fontSize['6xl'], fontWeight: fontWeight.extrabold, lineHeight: Math.round(fontSize['6xl'] * 1.2) },
  h2: { fontSize: fontSize['4xl'], fontWeight: fontWeight.bold, lineHeight: Math.round(fontSize['4xl'] * 1.2) },
  h3: { fontSize: fontSize['2xl'], fontWeight: fontWeight.bold, lineHeight: Math.round(fontSize['2xl'] * 1.3) },
  h4: { fontSize: fontSize.xl, fontWeight: fontWeight.semibold, lineHeight: Math.round(fontSize.xl * 1.3) },
  body: { fontSize: fontSize.lg, fontWeight: fontWeight.regular, lineHeight: Math.round(fontSize.lg * 1.5) },
  bodySm: { fontSize: fontSize.md, fontWeight: fontWeight.regular, lineHeight: Math.round(fontSize.md * 1.5) },
  caption: { fontSize: fontSize.sm, fontWeight: fontWeight.medium, lineHeight: Math.round(fontSize.sm * 1.4) },
  rating: { fontSize: fontSize['8xl'], fontWeight: fontWeight.extrabold, lineHeight: Math.round(fontSize['8xl'] * 1.2) },
  badge: { fontSize: fontSize.sm, fontWeight: fontWeight.semibold, lineHeight: Math.round(fontSize.sm * 1.4) },
  tab: { fontSize: fontSize.xs, fontWeight: fontWeight.medium, lineHeight: Math.round(fontSize.xs * 1.4) },
  button: { fontSize: 17, fontWeight: fontWeight.bold, lineHeight: Math.round(17 * 1.3) },
  buttonSm: { fontSize: fontSize.md, fontWeight: fontWeight.semibold, lineHeight: Math.round(fontSize.md * 1.3) },
} as const;
