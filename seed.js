const { Client } = require('pg');
require('dotenv').config();
const bcrypt = require('bcryptjs');

const client = new Client({
  user: process.env.DB_USER || 'postgres',
  host: process.env.DB_HOST || 'localhost',
  database: process.env.DB_NAME || 'register_system',
  password: process.env.DB_PASSWORD || 'postgres',
  port: process.env.DB_PORT || 5432,
});

async function seed() {
  await client.connect();
  console.log('Connected to database');

  await client.query('TRUNCATE TABLE student_projects, projects, users, affiliations RESTART IDENTITY CASCADE;');

  const affResult = await client.query(`
    INSERT INTO affiliations (std_name, description) VALUES
    ($1, $2),
    ($3, $4),
    ($5, $6)
    RETURNING id, std_name;
  `, [
    'Ban Chi dao Chien dich Mua He Xanh', 'Central指挥Committee',
    'Dai hoc Bach Khoa TP.HCM', 'University',
    'Dai hoc Kinh te - Luat', 'University'
  ]);

  const affMap = {};
  affResult.rows.forEach(r => { affMap[r.std_name] = r.id; });

  console.log('Affiliations:');
  affResult.rows.forEach(r => console.log(`  [${r.id}] ${r.std_name}`));

  const users = [
    {
      username: 'admin_root',
      full_name: 'Quan Tri Vien Tong',
      password_plaintext: 'admin123',
      role: 'admin',
      student_id: 'HQ-001',
      email: 'admin@mhx_system.vn',
      affiliation_id: affMap['Ban Chi dao Chien dich Mua He Xanh'],
    },
    {
      username: 'admin_deputy',
      full_name: 'Pho Quan Tri Vien',
      password_plaintext: 'admin456',
      role: 'admin',
      student_id: 'HQ-002',
      email: 'deputy@mhx_system.vn',
      affiliation_id: affMap['Ban Chi dao Chien dich Mua He Xanh'],
    },
    {
      username: 'cb_bachkhoa',
      full_name: 'Nguyen Van Can Bo',
      password_plaintext: 'school456',
      role: 'school',
      student_id: 'CB-BK-01',
      email: 'canbo@hcmut.edu.vn',
      affiliation_id: affMap['Dai hoc Bach Khoa TP.HCM'],
    },
    {
      username: 'cb_kinhteluat',
      full_name: 'Tran Thi Can Bo',
      password_plaintext: 'school789',
      role: 'school',
      student_id: 'CB-KTL-01',
      email: 'canbo@uel.edu.vn',
      affiliation_id: affMap['Dai hoc Kinh te - Luat'],
    },
    {
      username: 'leader_binh_phuoc',
      full_name: 'Tran Van Dia Phuong',
      password_plaintext: 'local789',
      role: 'community',
      student_id: 'LOC-BP-01',
      email: 'leader@binhphuoc.gov.vn',
      affiliation_id: affMap['Ban Chi dao Chien dich Mua He Xanh'],
    },
    {
      username: 'leader_daklak',
      full_name: 'Le Thi Dia Phuong',
      password_plaintext: 'local321',
      role: 'community',
      student_id: 'LOC-DL-01',
      email: 'leader@daklak.gov.vn',
      affiliation_id: affMap['Ban Chi dao Chien dich Mua He Xanh'],
    },
    {
      username: 'student_nhu',
      full_name: 'Nguyen Quynh Nhu',
      password_plaintext: 'student_nhu',
      role: 'student',
      student_id: 'SV2026_BK01',
      email: 'nhu.nguyen@student.hcmut.edu.vn',
      affiliation_id: affMap['Dai hoc Bach Khoa TP.HCM'],
    },
    {
      username: 'student_minh',
      full_name: 'Pham Duc Minh',
      password_plaintext: 'student_minh',
      role: 'student',
      student_id: 'SV2026_KTL01',
      email: 'minh.pham@student.uel.edu.vn',
      affiliation_id: affMap['Dai hoc Kinh te - Luat'],
    },
  ];

  const userMap = {};
  console.log('Users:');
  for (const u of users) {
    const password_hash = await bcrypt.hash(u.password_plaintext, 12);
    const res = await client.query(
      `INSERT INTO users (username, full_name, password_hash, role, is_active, student_id, email, affiliation_id)
       VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
      [u.username, u.full_name, password_hash, u.role, true, u.student_id, u.email, u.affiliation_id]
    );
    userMap[u.username] = res.rows[0].id;
    console.log(`  id=${res.rows[0].id} | ${u.username} | ${u.role} | password=${u.password_plaintext}`);
  }

  const affBachKhoa = affMap['Dai hoc Bach Khoa TP.HCM'];
  const affKinhTeLuat = affMap['Dai hoc Kinh te - Luat'];
  const leader_bp = userMap['leader_binh_phuoc'];
  const leader_dl = userMap['leader_daklak'];

  const projectResult = await client.query(`
    INSERT INTO projects (affiliation_id, community_user_id, name, description, num_max, project_start_day, project_end_day, form_start_day, form_end_day, date_approved, created_at)
    VALUES
      ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW()),
      ($11, $12, $13, $14, $15, $16, $17, $18, $19, $20, NOW()),
      ($21, $22, $23, $24, $25, $26, $27, $28, $29, $30, NOW())
    RETURNING id, name;
  `, [
    // leader_binh_phuoc creates 2 projects: 1 for Bach Khoa (approved), 1 for Kinh Te Luat (pending)
    affBachKhoa, leader_bp,
    'Mua He Xanh 2026 - Bach Khoa',
    'Tinh nguyen tai Dai hoc Bach Khoa TP.HCM',
    30,
    '2026-07-01T00:00:00Z', '2026-07-31T23:59:59Z',
    '2026-05-01T00:00:00Z', '2026-06-15T23:59:59Z',
    '2026-04-10T10:00:00Z',

    affKinhTeLuat, leader_bp,
    'Mua He Xanh 2026 - Kinh Te Luat',
    'Tinh nguyen tai Dai hoc Kinh te - Luat',
    25,
    '2026-07-01T00:00:00Z', '2026-07-31T23:59:59Z',
    '2026-05-01T00:00:00Z', '2026-06-15T23:59:59Z',
    null,

    // leader_daklak creates 1 project for Kinh Te Luat (approved)
    affKinhTeLuat, leader_dl,
    'Mua He Xanh 2026 - Kinh Te Luat Phase 2',
    'Dot 2 tinh nguyen tai Dai hoc Kinh te - Luat',
    20,
    '2026-08-01T00:00:00Z', '2026-08-31T23:59:59Z',
    '2026-06-16T00:00:00Z', '2026-07-15T23:59:59Z',
    '2026-04-12T10:00:00Z',
  ]);

  console.log('Projects:');
  projectResult.rows.forEach(r => console.log(`  [${r.id}] ${r.name}`));

  const approvedBachKhoa = projectResult.rows.find(r => r.name === 'Mua He Xanh 2026 - Bach Khoa').id;
  const approvedKinhTeLuat = projectResult.rows.find(r => r.name === 'Mua He Xanh 2026 - Kinh Te Luat Phase 2').id;

  const appResult = await client.query(`
    INSERT INTO student_projects (user_id, project_id, status, created_at)
    VALUES
      ($1, $2, $3, NOW()),
      ($4, $5, $6, NOW())
    RETURNING id, status;
  `, [
    userMap['student_nhu'], approvedBachKhoa, 'SCHOOL_PENDING',
    userMap['student_minh'], approvedKinhTeLuat, 'SCHOOL_PENDING',
  ]);

  console.log('Applications:');
  appResult.rows.forEach(r => console.log(`  id=${r.id} | status=${r.status}`));

  console.log(`\nSeeded ${users.length} users, ${projectResult.rows.length} projects, ${appResult.rows.length} applications`);
  await client.end();
}

seed().catch(err => { console.error(err.message); process.exit(1); });
