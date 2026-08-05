import { paths } from 'src/routes/paths';

// ----------------------------------------------------------------------

// Tambahkan menu baru SETELAH halamannya ada (path dari paths.ts, bukan '#').
export const navData = [
  { title: 'Home', path: paths.home },
  { title: 'Article', path: paths.article.root },
  { title: 'FAQ', path: paths.support },
];
