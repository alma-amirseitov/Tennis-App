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

export const lineHeight = {
  tight: 1.2,
  normal: 1.4,
  relaxed: 1.6,
} as const;

export const textStyles = {
  h1: { fontSize: fontSize['6xl'], fontWeight: fontWeight.extrabold, lineHeight: lineHeight.tight },
  h2: { fontSize: fontSize['4xl'], fontWeight: fontWeight.bold, lineHeight: lineHeight.tight },
  h3: { fontSize: fontSize['2xl'], fontWeight: fontWeight.bold, lineHeight: 1.3 },
  h4: { fontSize: fontSize.xl, fontWeight: fontWeight.semibold, lineHeight: 1.3 },
  body: { fontSize: fontSize.lg, fontWeight: fontWeight.regular, lineHeight: 1.5 },
  bodySm: { fontSize: fontSize.md, fontWeight: fontWeight.regular, lineHeight: 1.5 },
  caption: { fontSize: fontSize.sm, fontWeight: fontWeight.medium, lineHeight: 1.4 },
  rating: { fontSize: fontSize['8xl'], fontWeight: fontWeight.extrabold },
  badge: { fontSize: fontSize.sm, fontWeight: fontWeight.semibold },
  tab: { fontSize: fontSize.xs, fontWeight: fontWeight.medium },
  button: { fontSize: 17, fontWeight: fontWeight.bold },
  buttonSm: { fontSize: fontSize.md, fontWeight: fontWeight.semibold },
} as const;
