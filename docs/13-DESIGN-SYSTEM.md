# Design System
## Токены, компоненты и правила UI

---

## 1. Colors

### Primary Palette
```typescript
export const colors = {
  // Primary — Green (tennis court)
  primary:      '#0F8B4F',  // Main actions, active states
  primaryLight: '#E8F5EE',  // Backgrounds, highlights
  primaryDark:  '#0A6B3C',  // Pressed states, gradients

  // Accent — Gold (championship)
  accent:       '#F5A623',  // Secondary actions, highlights
  accentLight:  '#FEF3C7',  // Warning backgrounds

  // Neutrals
  background:   '#FAFBFC',  // App background
  card:         '#FFFFFF',  // Card background
  text:         '#1A1D21',  // Primary text
  textSecondary:'#6B7280',  // Secondary text, labels
  textMuted:    '#9CA3AF',  // Hints, timestamps, disabled
  border:       '#E5E7EB',  // Borders, dividers
  borderLight:  '#F3F4F6',  // Subtle borders, backgrounds

  // Semantic
  success:      '#10B981',  // Won, confirmed, online
  warning:      '#F59E0B',  // Pending, filling
  danger:       '#EF4444',  // Error, delete, lost
  info:         '#3B82F6',  // Info badges, links

  // Status
  statusOpen:   '#22C55E',  // Event open
  statusFilling:'#EAB308',  // Event filling
  statusFull:   '#3B82F6',  // Event full
  statusDone:   '#9CA3AF',  // Completed
  statusCancel: '#EF4444',  // Cancelled
};
```

### Usage Rules
- Primary green for main CTAs, active tab, links
- Never use raw hex in components — always reference `colors.X`
- Dark text on light backgrounds, light text on dark/primary backgrounds
- Status colors ONLY for status indicators, not general UI
- Badges use `color + '15'` as background (15% opacity)

---

## 2. Typography

### Font Family
```typescript
// Mobile: system fonts
export const fonts = {
  regular:  'System',        // -apple-system on iOS, Roboto on Android
  medium:   'System',        // weight 500
  semibold: 'System',        // weight 600
  bold:     'System',        // weight 700
  extrabold:'System',        // weight 800
};

// Web: Inter (via Google Fonts) + system fallback
// font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
```

### Type Scale
```typescript
export const fontSize = {
  xs:   11,   // Tiny labels, badge counts
  sm:   12,   // Captions, metadata, timestamps
  base: 13,   // Secondary text, descriptions
  md:   14,   // Body text, list items
  lg:   15,   // Card titles, input text
  xl:   16,   // Section headers
  '2xl':18,   // Screen titles (secondary)
  '3xl':20,   // Large headers
  '4xl':22,   // Screen title
  '5xl':24,   // Hero numbers
  '6xl':28,   // Large hero
  '7xl':32,   // XL hero
  '8xl':36,   // Rating display
};

export const fontWeight = {
  regular:   '400',
  medium:    '500',
  semibold:  '600',
  bold:      '700',
  extrabold: '800',
};

export const lineHeight = {
  tight:  1.2,  // Headlines
  normal: 1.4,  // Body text
  relaxed:1.6,  // Long-form text
};
```

### Preset Styles
```typescript
export const textStyles = {
  // Headlines
  h1: { fontSize: 28, fontWeight: '800', lineHeight: 1.2 },  // Screen title
  h2: { fontSize: 22, fontWeight: '700', lineHeight: 1.2 },  // Section title
  h3: { fontSize: 18, fontWeight: '700', lineHeight: 1.3 },  // Card title
  h4: { fontSize: 16, fontWeight: '600', lineHeight: 1.3 },  // Subsection

  // Body
  body:    { fontSize: 15, fontWeight: '400', lineHeight: 1.5 },
  bodySm:  { fontSize: 14, fontWeight: '400', lineHeight: 1.5 },
  caption: { fontSize: 12, fontWeight: '500', lineHeight: 1.4 },

  // Special
  rating:  { fontSize: 36, fontWeight: '800' },   // Big rating number
  badge:   { fontSize: 12, fontWeight: '600' },    // Badge text
  tab:     { fontSize: 10, fontWeight: '500' },    // Tab bar label
  button:  { fontSize: 17, fontWeight: '700' },    // Primary button
  buttonSm:{ fontSize: 14, fontWeight: '600' },    // Small button
};
```

---

## 3. Spacing

```typescript
export const spacing = {
  xs:  4,
  sm:  8,
  md:  12,
  lg:  16,
  xl:  20,
  '2xl': 24,
  '3xl': 32,
  '4xl': 40,
  '5xl': 48,
  '6xl': 64,
};

// Screen padding: 16 (lg)
// Card padding: 16 (lg)
// Card gap: 12 (md) or 16 (lg)
// Section gap: 24 (2xl)
// List item padding: 14 vertical, 16 horizontal
// Input padding: 14 vertical, 16 horizontal
```

---

## 4. Border Radius

```typescript
export const radius = {
  sm:   8,    // Small buttons, chips
  md:   12,   // Inputs, small cards
  lg:   16,   // Cards, modals
  xl:   20,   // Large cards
  pill: 100,  // Badges, avatar
  full: 9999, // Circle
};
```

---

## 5. Shadows

```typescript
export const shadows = {
  sm: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 2,
    elevation: 1,
  },
  md: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.08,
    shadowRadius: 8,
    elevation: 3,
  },
  lg: {
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.12,
    shadowRadius: 16,
    elevation: 5,
  },
};
```

---

## 6. Component Specs

### Button
```
Primary:  bg=primary, text=white, h=52, radius=14, font=button
Secondary: bg=primaryLight, text=primary, border=primary, h=52, radius=14
Outline:  bg=transparent, text=text, border=border, h=52, radius=14
Small:    h=36, radius=8, font=buttonSm, padding=6,16
Disabled: opacity=0.5, no press
Loading:  spinner replaces text, disabled
```

### Input
```
Default:  bg=background, border=border, h=52, radius=12, padding=14,16
Focused:  border=primary (2px)
Error:    border=danger, error text below in danger color
Label:    above input, fontSize=13, fontWeight=600, color=textSecondary
```

### Card
```
bg=card, radius=16, padding=16, border=borderLight (1px)
Interactive: cursor=pointer (web), Pressable (mobile)
No shadow by default (flat design)
```

### Avatar
```
Sizes: 24 (xs), 32 (sm), 40 (md), 48 (lg), 64 (xl), 80 (xxl)
Shape: circle (borderRadius=size/2)
Fallback: initials on colored background (primary+20% opacity)
Online indicator: 12px green circle, bottom-right, white border
```

### Badge
```
padding: 3,10
radius: pill (20)
bg: color + 15% opacity
text: color, fontSize=12, fontWeight=600
Variants: primary, success, warning, danger, info, muted
```

### Tab Bar
```
Height: 80px (includes safe area)
Background: card
Border top: borderLight (1px)
Icon: 22px, active=full opacity, inactive=0.4
Label: 10px, active=primary+bold, inactive=muted
Items: 5 (Home, Players, Events, Communities, Profile)
```

### Header
```
Height: 52px
Background: card
Border bottom: borderLight (1px)
Title: fontSize=18, fontWeight=700
Back button: ← (fontSize=22)
Right icons: 20px with notification badges
```

### Skeleton Loading
```
Background: borderLight (#F3F4F6)
Animation: pulse (opacity 0.4 → 1.0 → 0.4, 1.5s loop)
Shape matches content shape (text=rounded rect, avatar=circle)
```

### Empty State
```
Center of screen
Large icon/emoji: 48px
Title: h3, marginTop=16
Description: bodySm, color=textMuted, marginTop=4
Action button: Primary or Secondary, marginTop=16
```

### Error State
```
Center of screen
Icon: ⚠️ or relevant emoji
Title: "Что-то пошло не так"
Description: error message or generic
Retry button: Primary
```

---

## 7. Iconography

### Approach
Emoji-first for MVP (no icon library dependency). Consistent emoji set:

```
Navigation: 🏠 👥 🎾 🏛 👤
Actions:    ➕ ✏️ 🗑 ↩️ 📤
Status:     ● (colored circles for status dots)
Social:     ❤️ 💬 🔔 👀
Sport:      🎾 🏆 🏅 🎖 ⭐ 🔥
Profile:    📊 📅 ⚙️ 📍 🎯
Chat:       💬 ✓✓ (read receipts)
```

### Post-MVP
Migrate to Lucide icons (react-native-lucide) for more polished look.

---

## 8. Animation Guidelines

### Principles
- Subtle, not distracting
- 200-300ms duration for most transitions
- Ease-out for elements entering, ease-in for leaving

### Specific Animations
```
Screen transitions: slide from right (300ms)
Tab switch: cross-fade (200ms)
Button press: scale 0.97 (100ms)
Card press: scale 0.98 (100ms)
List item appear: fade-in + slide-up (staggered, 50ms delay each)
Pull-to-refresh: native (ScrollView)
Skeleton pulse: opacity 0.4→1.0 (1.5s, infinite)
Badge count: scale bounce (spring)
Toast: slide from top (300ms) + auto-dismiss (3s)
```

---

## 9. Responsive (Web)

### Breakpoints
```typescript
// Web admin only (mobile app is mobile-only)
export const breakpoints = {
  sm:  640,   // Mobile
  md:  768,   // Tablet
  lg:  1024,  // Desktop
  xl:  1280,  // Wide desktop
};
```

### Layout
```
Sidebar: 256px fixed (desktop), drawer (tablet/mobile)
Content: fluid, max-width 1200px, centered
Tables: horizontal scroll on mobile
Cards: 1 column mobile, 2 tablet, 3-4 desktop
```

---

## 10. Dark Mode

Not included in MVP. Prepare by:
- Using `colors` tokens everywhere (never hardcode hex)
- Using semantic names (background, card, text) not literal names (white, black)
- Future: swap color tokens for dark palette
