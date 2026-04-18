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

  console.log(`\nSeeded ${users.length} users, ${projectResult.rows.length} projects, ${appResult.rows.length} applications`);
  await client.end();
}

seed().catch(err => { console.error(err.message); process.exit(1); });
