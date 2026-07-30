'use client';

import type { ArticleCategory } from 'src/lib/api';

import { debounce } from 'es-toolkit';
import { useRef, useMemo, useState, useEffect } from 'react';

import Box from '@mui/material/Box';
import Chip from '@mui/material/Chip';
import Skeleton from '@mui/material/Skeleton';
import TextField from '@mui/material/TextField';
import InputAdornment from '@mui/material/InputAdornment';

import { Iconify } from 'src/components/iconify';

// ----------------------------------------------------------------------

const SEARCH_DEBOUNCE_MS = 400;

type ArticleFiltersProps = {
  search: string;
  category: string;
  categories: ArticleCategory[];
  loadingCategories?: boolean;
  onSearchChange: (value: string) => void;
  onCategoryChange: (slug: string) => void;
};

export function ArticleFilters({
  search,
  category,
  categories,
  loadingCategories = false,
  onSearchChange,
  onCategoryChange,
}: ArticleFiltersProps) {
  // Controlled input synced FROM the URL only on external changes (browser
  // back/forward) — a naive sync would clobber keystrokes typed while the
  // debounced navigation is still in flight.
  const [value, setValue] = useState(search);
  const lastEmitted = useRef(search);

  // onSearchChange's identity rotates on every committed navigation (it
  // closes over useSearchParams()). Route it through a ref so the debounced
  // function stays stable — recreating it would cancel pending trailing
  // keystrokes mid-navigation.
  const onSearchChangeRef = useRef(onSearchChange);
  useEffect(() => {
    onSearchChangeRef.current = onSearchChange;
  }, [onSearchChange]);

  const debouncedSearchChange = useMemo(
    () =>
      debounce((nextValue: string) => {
        lastEmitted.current = nextValue;
        onSearchChangeRef.current(nextValue);
      }, SEARCH_DEBOUNCE_MS),
    []
  );

  useEffect(() => () => debouncedSearchChange.cancel(), [debouncedSearchChange]);

  useEffect(() => {
    if (search !== lastEmitted.current) {
      lastEmitted.current = search;
      setValue(search);
    }
  }, [search]);

  return (
    <Box
      sx={{
        gap: 3,
        display: 'flex',
        flexDirection: { xs: 'column', md: 'row' },
        alignItems: { xs: 'flex-start', md: 'center' },
        justifyContent: 'space-between',
      }}
    >
      <Box sx={{ gap: 1, display: 'flex', flexWrap: 'wrap' }}>
        <Chip
          label="All"
          variant={category ? 'outlined' : 'filled'}
          onClick={() => onCategoryChange('')}
        />

        {loadingCategories
          ? Array.from({ length: 3 }, (_, index) => (
              <Skeleton key={index} variant="rounded" sx={{ width: 72, height: 32 }} />
            ))
          : categories.map((item) => (
              <Chip
                key={item.slug}
                label={item.name}
                variant={category === item.slug ? 'filled' : 'outlined'}
                onClick={() => onCategoryChange(category === item.slug ? '' : item.slug)}
              />
            ))}
      </Box>

      <TextField
        fullWidth
        size="small"
        hiddenLabel
        value={value}
        placeholder="Search articles..."
        onChange={(event) => {
          setValue(event.target.value);
          debouncedSearchChange(event.target.value);
        }}
        slotProps={{
          input: {
            startAdornment: (
              <InputAdornment position="start">
                <Iconify icon="carbon:search" sx={{ color: 'text.disabled' }} />
              </InputAdornment>
            ),
          },
        }}
        sx={{ maxWidth: { md: 320 } }}
      />
    </Box>
  );
}
