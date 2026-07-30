import { CONFIG } from 'src/global-config';

// ----------------------------------------------------------------------

export const asset = (name: string) => `${CONFIG.assetsDir}/assets/venturo/${name}`;

export const CONTACT = {
  wa: 'https://wa.me/6285128043814',
  email: 'hello@venturo.id',
};

// ----------------------------------------------------------------------

export const HERO = {
  badge: 'Konsultasi 100% Gratis',
  title: 'Software House Malang dengan Talenta Programmer Terbesar',
  description:
    'Venturo adalah jasa outsource programmer terbesar di Malang, didukung 130+ talenta programmer, dan dipercaya berbagai bisnis serta brand untuk layanan pembuatan software.',
  cta: 'Mulai Konsultasi',
  stats: [
    { value: 130, suffix: '+', label: 'Dedicated IT experts' },
    { value: 100, suffix: '+', label: 'Successful Projects' },
  ],
  trustedLabel: 'Dipercaya Oleh:',
};

/**
 * Logo klien/brand untuk row "Dipercaya Oleh" (trusted-by).
 * File logo ada di `public/assets/venturo/clients/`.
 * Untuk menambah klien: taruh file di folder itu lalu daftarkan di sini.
 */
export const CLIENTS: { name: string; logo: string }[] = [
  { name: 'Qoin', logo: asset('clients/qoin.webp') },
  { name: 'Shipper', logo: asset('clients/shipper.webp') },
  { name: 'Hayyu Skin Clinic', logo: asset('clients/hayyu.webp') },
  { name: 'Bobobox', logo: asset('clients/bobobox.webp') },
  { name: 'Liputan 6', logo: asset('clients/liputan6.webp') },
  { name: 'Powmeals', logo: asset('clients/powmeals.webp') },
  { name: 'Bolong', logo: asset('clients/bolong.webp') },
  { name: 'Majoo', logo: asset('clients/majoo.webp') },
  { name: 'Pudu', logo: asset('clients/pudu.webp') },
  { name: 'Humanis', logo: asset('clients/humanis.webp') },
  { name: 'Wisdom Crowd', logo: asset('clients/wisdom-crowd.webp') },
  { name: 'Bio Farma', logo: asset('clients/biofarma.webp') },
  { name: 'Kemenkes', logo: asset('clients/kemenkes.webp') },
  { name: 'Olahkarsa', logo: asset('clients/olahkarsa.webp') },
];

// ----------------------------------------------------------------------

export const PROBLEM = {
  caption: 'The Problem',
  title: 'Proyek IT Anda Bermasalah?',
  // TODO: user akan kirim foto "stressed man" → taruh di public/assets/venturo/problem-stress.webp
  image: asset('problem-stress.webp'),
  items: [
    {
      title: 'Kesulitan Menemukan Programmer yang Tepat',
      description:
        'Lulusan IT melimpah, tapi menemukan programmer andal sesuai harapan bukanlah hal yang mudah.',
      icon: 'solar:user-rounded-outline',
    },
    {
      title: 'Proyek Molor, Melebihi Anggaran, Tak Sesuai Ekspektasi',
      description:
        'Efeknya proyek makin lama selesai, budget makin bengkak, Anda makin kecewa. Anda tidak sendirian — banyak bisnis juga mengalaminya.',
      icon: 'solar:hourglass-line-outline',
    },
    {
      title: 'Tingkat Perputaran yang Tinggi di Industri Ini',
      description:
        'Bosan gonta-ganti programmer! Tingkat turnover yang relatif tinggi mencapai 19,22% — lebih tinggi daripada industri lain.',
      icon: 'solar:restart-bold',
    },
  ],
  closing: 'Dan akhirnya… bisnis Anda kehilangan waktu, biaya, bahkan kepercayaan klien.',
};

// ----------------------------------------------------------------------

export const SOLUTION = {
  caption: 'Stress-Free Solutions',
  title: 'Solusi Tanpa Drama untuk Proyek IT Anda',
  image: asset('tim-expert-programmers.webp'),
  items: [
    {
      title: 'Tim yang Berdedikasi',
      icon: asset('ikon-tim-programmer-berdedikasi.webp'),
      description:
        'Bangun tim outsource full-time layaknya memiliki tim internal sendiri. Squad ini sepenuhnya fokus mengerjakan proyek Anda tanpa terbagi ke proyek lain.',
    },
    {
      title: 'Supervisor Berpengalaman',
      icon: asset('ikon-quality-assurance-supervisi.webp'),
      description:
        'Setiap tim didukung supervisor senior yang memastikan kualitas melalui analisa, pendampingan, dan review dari berbagai proyek serupa.',
    },
    {
      title: 'Respon Cepat',
      icon: asset('ikon-respon-cepat.webp'),
      description:
        'Selama masa kontrak berjalan, tim kami siap memberikan dukungan respons cepat melalui grup messenger online bersama seluruh PIC proyek Anda.',
    },
    {
      title: 'Tes Otomatisasi',
      icon: asset('ikon-otomatisasi-pengembangan-software.webp'),
      description:
        'Pengujian otomatis untuk memastikan kualitas perangkat lunak secara cepat, akurat, konsisten, dan meminimalisir human error.',
    },
    {
      title: 'Laporan Progresif',
      icon: asset('ikon-laporan-progres-proyek.webp'),
      description:
        'Laporan progress detail diberikan setiap minggu, sehingga Anda selalu mendapat update perkembangan proyek secara transparan.',
    },
  ],
};

// ----------------------------------------------------------------------

export const FOCUS = {
  caption: 'Focus on Your Core Business',
  title: 'Fokus pada Bisnis Inti Anda, Biar Urusan IT Kami yang Pegang',
  description:
    'Venturo, jasa pembuatan software house Malang, siap membantu dengan programmer profesional yang bekerja seolah tim in-house Anda sendiri.',
  highlight: 'Tim IT Khusus Anda',
  image: asset('diagram-outsourcing-tim-it-venturo.webp'),
  roles: [
    'Project Manager',
    'System Analyst',
    'UI/UX Designer',
    'Programmer / Developer',
    'Quality Assurance (QA)',
    'Automation Test Engineer',
    'Technical Documentation',
  ],
};

// ----------------------------------------------------------------------

export const MANAGEMENT = {
  caption: 'We Already Have a Good Management',
  title: 'Manajemen Proyek yang Sudah Teruji',
  description: 'Jasa outsourcing programmer, perusahaan IT di Malang terbaik.',
  phone: asset('app-phone.webp'),
  items: [
    {
      title: 'Roadmap untuk Timeline Planning',
      icon: asset('ikon-penjadwalan-proyek.webp'),
      description:
        'Roadmap berisi list task dan timeline project yang dapat Anda pantau setiap saat, dibuat berdasarkan kesepakatan dengan Anda.',
    },
    {
      title: 'Monitoring Sprint Mingguan',
      icon: asset('ikon-monitoring-proyek.webp'),
      description:
        'Mengontrol sprint dan memastikan programmer kami bekerja sesuai dengan waktu yang Anda beli.',
    },
    {
      title: 'Lacak Progress Terkini',
      icon: asset('ikon-laporan-progres-proyek.webp'),
      description:
        'Dengan sekali lihat, Anda dapat memantau modul mana yang sedang berjalan atau sudah selesai.',
    },
    {
      title: 'Tes Skenario untuk Mencegah Bug',
      icon: asset('ikon-otomatisasi-pengembangan-software.webp'),
      description:
        'Setiap use case dicatat secara rinci pada Test Scenario sebagai langkah antisipasi terhadap bug/error.',
    },
    {
      title: 'Layanan Maksimal untuk Klien',
      icon: asset('garansi-bug-project-selesai.webp'),
      description:
        'Kami menjamin pelayanan maksimal 48 jam penanganan bug dan improvement yang dapat Anda pantau.',
    },
    {
      title: 'Memantau Kinerja Setiap Squad',
      icon: asset('ikon-quality-assurance-supervisi.webp'),
      description:
        'Setiap talenta memiliki target harian dan bulanan, dan bersaing untuk menjadi yang terbaik.',
    },
  ],
};

// ----------------------------------------------------------------------

export const RESOURCE = {
  titleLines: ['The Largest', 'Programmer Resource', 'in Malang'],
  description:
    'Cari programmer terbaik di Malang? Venturo menghadirkan the largest programmer resource in Malang untuk kebutuhan outsourcing dan pengembangan software skala besar.',
  videoUrl: 'https://www.youtube.com/watch?v=1W35KcCQqww',
  videoThumb: asset('video-thumb-manajemen-proyek.webp'),
  videoLabel: 'Supervisor & Project Manager',
};

// ----------------------------------------------------------------------

export const SPECIAL_OFFER = {
  heading: 'SPECIAL OFFER!!',
  promoStrong: '50% OFF hiring fee',
  promo: 'khusus hingga 31 Mei 2026, siap menyelesaikan proyek IT Anda tanpa drama!',
  cta: 'Hubungi Kami Sekarang',
  note: 'Konsultasi gratis untuk kebutuhan bisnis Anda',
  image: asset('person-arms.webp'),
};

// ----------------------------------------------------------------------

export const TECH_STACK = {
  caption: 'Technology Stack',
  title: 'Teknologi yang Kami Kuasai',
  groups: [
    {
      label: 'Web Programmer',
      logos: [
        { name: 'Laravel', logo: 'logo-laravel.webp' },
        { name: 'Go', logo: 'logo-golang.webp' },
        { name: 'Node.js', logo: 'logo-nodejs.webp' },
        { name: 'React', logo: 'logo-react.webp' },
        { name: 'Vue.js', logo: 'logo-vuejs.webp' },
        { name: 'Angular', logo: 'logo-angularjs.webp' },
        { name: 'Bootstrap', logo: 'logo-bootstrap.webp' },
        { name: 'Tailwind CSS', logo: 'logo-tailwind-css.webp' },
        { name: 'MySQL', logo: 'logo-mysql.webp' },
        { name: 'MongoDB', logo: 'logo-mongodb.webp' },
        { name: 'Redis', logo: 'logo-redis.webp' },
        { name: 'RabbitMQ', logo: 'logo-rabbitmq.webp' },
        { name: 'Sentry', logo: 'logo-sentry.webp' },
      ],
    },
    {
      label: 'Mobile Programmer',
      logos: [
        { name: 'Flutter', logo: 'logo-flutter.webp' },
        { name: 'Kotlin', logo: 'logo-kotlin.webp' },
        { name: 'Swift', logo: 'logo-swift.webp' },
        { name: 'Firebase', logo: 'logo-firebase.webp' },
      ],
    },
    {
      label: 'UI / UX',
      logos: [{ name: 'Figma', logo: 'logo-figma.webp' }],
    },
  ],
};

// ----------------------------------------------------------------------

export const FAQS = {
  caption: 'FAQ',
  title: 'Frequently Asked Questions',
  items: [
    {
      question:
        'Apa keunggulan menggunakan layanan outsourcing programmer dibanding merekrut langsung?',
      answer:
        'Anda mendapatkan tim siap kerja tanpa repot rekrutmen, training, dan retensi. Hemat waktu dan biaya, dengan supervisor serta manajemen yang sudah terbukti, sehingga Anda bisa fokus ke core business.',
    },
    {
      question: 'Apakah tim Anda berpengalaman dengan proyek AI atau machine learning?',
      answer:
        'Ya. Tim kami menangani berbagai proyek termasuk integrasi AI/ML — mulai dari pengolahan data, model prediktif, hingga implementasi fitur berbasis AI pada aplikasi.',
    },
    {
      question: 'Bagaimana model kerja outsourcing programmer di Venturo?',
      answer:
        'Anda membentuk dedicated squad full-time yang bekerja seperti tim in-house Anda, didampingi supervisor, dengan roadmap, sprint mingguan, dan laporan progres transparan.',
    },
    {
      question: 'Berapa lama waktu Venturo untuk mempersiapkan tim hingga siap untuk Development?',
      answer:
        'Umumnya tim dapat disiapkan dalam hitungan hari setelah kebutuhan dan kesepakatan final, tergantung skala serta stack yang dibutuhkan.',
    },
    {
      question: 'Apakah saya menerima source code setelah proses development selesai?',
      answer: 'Ya. Seluruh source code menjadi milik Anda dan diserahkan sesuai kesepakatan kerja.',
    },
    {
      question: 'Apakah ada garansi setelah melakukan development?',
      answer:
        'Ada. Kami pastikan sistem bebas dari bug kritis pasca serah terima, dengan garansi perbaikan hingga 30 hari tanpa biaya tambahan.',
    },
    {
      question: 'Bagaimana jika programmer yang diberikan tidak cocok?',
      answer:
        'Kami akan melakukan evaluasi dan penggantian talenta agar tim benar-benar sesuai dengan kebutuhan dan ekspektasi proyek Anda.',
    },
  ],
};

// ----------------------------------------------------------------------

export const CLOSING_CTA = {
  title: 'Garansi Bug Setelah Project Selesai',
  descriptionStrong: '30 Hari Tanpa Biaya Tambahan',
  description:
    'Kami pastikan sistem bebas dari bug kritis pasca serah terima, dengan garansi perbaikan hingga',
  descriptionEnd: 'Jika Anda serius maka tim kami siap mendukung dari awal hingga selesai.',
  cta: 'Mulai Konsultasi',
  image: asset('person-laptop.webp'),
};
