import { View, Text, Image, StyleSheet } from 'react-native';
import { colors, radius } from '@/shared/theme';

export type AvatarSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl';

const SIZE_MAP: Record<AvatarSize, number> = {
  xs: 24,
  sm: 32,
  md: 40,
  lg: 56,
  xl: 80,
};

export interface AvatarProps {
  uri?: string | null;
  name?: string | null;
  size?: AvatarSize;
  showOnline?: boolean;
}

function getInitials(name: string | undefined | null): string {
  if (!name || name.trim() === '') return '?';
  const parts = name.trim().split(/\s+/);
  if (parts.length >= 2) {
    const first = parts[0]?.[0] ?? '';
    const last = parts[parts.length - 1]?.[0] ?? '';
    return `${first}${last}`.toUpperCase();
  }
  return (parts[0]?.[0] ?? '?').toUpperCase();
}

export function Avatar({
  uri,
  name,
  size = 'md',
  showOnline = false,
}: AvatarProps) {
  const s = SIZE_MAP[size];
  const halfSize = s / 2;

  return (
    <View style={[styles.wrapper, { width: s, height: s }]}>
      {uri ? (
        <Image
          source={{ uri }}
          style={[styles.image, { width: s, height: s, borderRadius: halfSize }]}
          resizeMode="cover"
        />
      ) : (
        <View
          style={[
            styles.fallback,
            { width: s, height: s, borderRadius: halfSize },
          ]}
        >
          <Text
            style={[
              styles.initials,
              { fontSize: s * 0.4 },
            ]}
          >
            {getInitials(name)}
          </Text>
        </View>
      )}
      {showOnline ? (
        <View
          style={[
            styles.onlineIndicator,
            {
              width: 12,
              height: 12,
              borderRadius: 6,
              right: 0,
              bottom: 0,
              borderWidth: 2,
              borderColor: colors.card,
            },
          ]}
        />
      ) : null}
    </View>
  );
}

const styles = StyleSheet.create({
  wrapper: {
    position: 'relative',
  },
  image: {},
  fallback: {
    backgroundColor: `${colors.primary}33`,
    alignItems: 'center',
    justifyContent: 'center',
  },
  initials: {
    fontWeight: '600',
    color: colors.primary,
  },
  onlineIndicator: {
    position: 'absolute',
    backgroundColor: colors.success,
  },
});
