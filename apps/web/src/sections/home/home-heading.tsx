import type { BoxProps } from '@mui/material/Box';

import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';

// ----------------------------------------------------------------------

type Props = BoxProps & {
  title: string;
  caption?: string;
  description?: string;
};

export function HomeHeading({ title, caption, description, sx, ...other }: Props) {
  return (
    <Box
      sx={[
        { mx: 'auto', maxWidth: 720, textAlign: 'center', mb: { xs: 5, md: 8 } },
        ...(Array.isArray(sx) ? sx : [sx]),
      ]}
      {...other}
    >
      {caption && (
        <Typography variant="overline" sx={{ mb: 2, display: 'block', color: 'primary.main' }}>
          {caption}
        </Typography>
      )}

      <Typography variant="h2">{title}</Typography>

      {description && (
        <Typography sx={{ mt: 2, color: 'text.secondary' }}>{description}</Typography>
      )}
    </Box>
  );
}
