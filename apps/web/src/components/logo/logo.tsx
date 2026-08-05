'use client';

import type { LinkProps } from '@mui/material/Link';

import { mergeClasses } from 'minimal-shared/utils';

import Box from '@mui/material/Box';
import Link from '@mui/material/Link';
import { styled } from '@mui/material/styles';

import { RouterLink } from 'src/routes/components';

import { CONFIG } from 'src/global-config';

import { logoClasses } from './classes';

// ----------------------------------------------------------------------

const fullLogoSrc = `${CONFIG.assetsDir}/assets/venturo/logo-venturo-software-house-malang.webp`;
const singleLogoSrc = `${CONFIG.assetsDir}/assets/venturo/logo-venturo-icon.webp`;

export type LogoProps = LinkProps & {
  isSingle?: boolean;
  disabled?: boolean;
};

export function Logo({
  sx,
  disabled,
  className,
  href = '/',
  isSingle = false,
  ...other
}: LogoProps) {
  return (
    <LogoRoot
      component={RouterLink}
      href={href}
      aria-label="Venturo"
      underline="none"
      className={mergeClasses([logoClasses.root, className])}
      sx={[
        {
          width: 'auto',
          height: 24,
          ...(isSingle && { width: 40, height: 40 }),
          ...(disabled && { pointerEvents: 'none' }),
        },
        ...(Array.isArray(sx) ? sx : [sx]),
      ]}
      {...other}
    >
      <Box
        component="img"
        alt="Venturo"
        src={isSingle ? singleLogoSrc : fullLogoSrc}
        sx={{ height: '100%', width: isSingle ? '100%' : 'auto', objectFit: 'contain' }}
      />
    </LogoRoot>
  );
}

// ----------------------------------------------------------------------

const LogoRoot = styled(Link)(() => ({
  flexShrink: 0,
  color: 'inherit',
  display: 'inline-flex',
  verticalAlign: 'middle',
}));
