import Box from '@mui/material/Box';
import Drawer from '@mui/material/Drawer';
import { styled } from '@mui/material/styles';
import Typography from '@mui/material/Typography';
import ButtonBase from '@mui/material/ButtonBase';
import Tab, { tabClasses } from '@mui/material/Tab';
import Tabs, { tabsClasses } from '@mui/material/Tabs';

import { Iconify } from 'src/components/iconify';
import { Scrollbar } from 'src/components/scrollbar';

// Data-only import: kontak fallback yang sama dengan section home.
import { CONTACT } from 'src/sections/home/home-data';

// ----------------------------------------------------------------------

type Props = {
  open: boolean;
  topic: string;
  onClose: () => void;
  waLink?: string | null;
  data: {
    title: string;
    icon: string;
  }[];
  onChangeTopic: (event: React.SyntheticEvent, newValue: string) => void;
};

export function SupportNav({ topic, data, onChangeTopic, open, onClose, waLink }: Props) {
  const renderItems = () => (
    <Tabs
      value={topic}
      onChange={onChangeTopic}
      orientation="vertical"
      sx={{ [`& .${tabsClasses.list}`]: { gap: 0 } }}
    >
      {data.map((item) => (
        <Tab
          key={item.title}
          value={item.title}
          label={item.title}
          icon={<Box component="img" alt="" src={item.icon} sx={{ width: 28, height: 28 }} />}
          sx={{
            gap: 1,
            typography: 'body2',
            justifyContent: 'flex-start',
            [`& .${tabClasses.selected}`]: {
              typography: 'subtitle2',
              fontWeight: 'fontWeightBold',
            },
          }}
        />
      ))}
    </Tabs>
  );

  const renderInfo = () => (
    <>
      <Typography component="h6" variant="h6" sx={{ mb: 1 }}>
        Masih butuh bantuan?
      </Typography>

      <Typography variant="body2" sx={{ color: 'text.secondary', mb: 3 }}>
        Tim kami siap menjawab pertanyaan Anda.
      </Typography>

      <SupportButton
        component="a"
        href={waLink ?? CONTACT.wa}
        target="_blank"
        rel="noopener noreferrer"
      >
        <Iconify width={24} icon="solar:chat-line-outline" />
        WhatsApp
      </SupportButton>

      <SupportButton component="a" href={`mailto:${CONTACT.email}`}>
        <Iconify width={24} icon="solar:letter-outline" />
        {CONTACT.email}
      </SupportButton>
    </>
  );

  return (
    <>
      <Box
        sx={(theme) => ({
          width: 280,
          flexShrink: 0,
          flexDirection: 'column',
          display: { xs: 'none', md: 'flex' },
          borderRight: `solid 1px ${theme.vars.palette.divider}`,
        })}
      >
        {renderItems()}
        <Box sx={{ mt: 3, pr: 5 }}>{renderInfo()}</Box>
      </Box>

      <Drawer
        open={open}
        onClose={onClose}
        slotProps={{
          paper: { sx: { width: 280 } },
        }}
      >
        <Scrollbar>
          <Box sx={{ pt: 2, pl: 2 }}>{renderItems()}</Box>

          <Box sx={{ p: 2.5 }}>{renderInfo()}</Box>
        </Scrollbar>
      </Drawer>
    </>
  );
}

// ----------------------------------------------------------------------

const SupportButton = styled(ButtonBase)(({ theme }) => ({
  ...theme.typography.subtitle2,
  width: '100%',
  alignItems: 'center',
  gap: theme.spacing(1.5),
  justifyContent: 'flex-start',
  padding: theme.spacing(1.5, 2),
  marginBottom: theme.spacing(2),
  borderRadius: theme.shape.borderRadius,
  border: `solid 1px ${theme.vars.palette.divider}`,
})) as typeof ButtonBase;
