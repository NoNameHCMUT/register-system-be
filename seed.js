const { Client } = require('pg');
require('dotenv').config(); // Tải biến môi trường từ file .env

// Cấu hình kết nối sử dụng các biến môi trường
const client = new Client({
  user: process.env.DB_USER,
  host: process.env.DB_HOST || 'localhost',
  database: process.env.DB_NAME,
  password: process.env.DB_PASSWORD,
  port: process.env.DB_PORT,
});

async function seedData() {
  try {
    await client.connect();
    console.log("--- Đã kết nối tới Database: " + process.env.DB_NAME + " ---");

    // 1. Dọn dẹp dữ liệu cũ để tránh trùng lặp khi chạy lại seed
    await client.query('TRUNCATE TABLE users, affiliations RESTART IDENTITY CASCADE;');

    // 2. Seed bảng affiliations (Tổ chức) dựa trên mô hình phối hợp ĐH và cộng đồng
    const affResult = await client.query(`
      INSERT INTO affiliations (std_name) VALUES 
      ('Ban Chỉ đạo Chiến dịch Mùa Hè Xanh'), 
      ('Đại học Bách Khoa TP.HCM'), 
      ('Đại học Kinh tế - Luật')
      RETURNING id, std_name;
    `);
    
    // Ánh xạ tên trường sang ID để gán cho User
    const affMap = {};
    affResult.rows.forEach(row => {
      affMap[row.std_name] = row.id;
    });

    console.log('>> Danh sách affiliations đã tạo:' );
    affResult.rows.forEach((row) => {
      console.log(`   - [${row.id}] ${row.std_name}`);
    });

    // 3. Danh sách người dùng mẫu (Dành cho 4 vai trò chính của hệ thống)
    const users = [
      {
        username: 'admin_root',
        full_name: 'Quản Trị Viên Tổng',
        password_hash: '$2a$12$fakehash_admin123',
        role: 'admin',
        student_id: 'HQ-001',
        email: 'admin@mhx_system.vn',
        affiliation_id: affMap['Ban Chỉ đạo Chiến dịch Mùa Hè Xanh']
      },
      {
        username: 'cb_bachkhoa',
        full_name: 'Nguyễn Văn Cán Bộ',
        password_hash: '$2a$12$fakehash_school456',
        role: 'school',
        student_id: 'CB-BK-01',
        email: 'canbo@hcmut.edu.vn',
        affiliation_id: affMap['Đại học Bách Khoa TP.HCM']
      },
      {
        username: 'leader_binh_phuoc',
        full_name: 'Trần Văn Địa Phương',
        password_hash: '$2a$12$fakehash_local789',
        role: 'community',
        student_id: 'LOC-BP-01',
        email: 'leader@binhphuoc.gov.vn',
        // schema requires affiliation_id NOT NULL; use the campaign affiliation for community leader
        affiliation_id: affMap['Ban Chỉ đạo Chiến dịch Mùa Hè Xanh']
      },
      {
        username: 'student_nhu',
        full_name: 'Nguyễn Quỳnh Như',
        password_hash: '$2a$12$fakehash_student_nhu',
        role: 'student',
        student_id: 'SV2026_BK01',
        email: 'nhu.nguyen@student.hcmut.edu.vn',
        affiliation_id: affMap['Đại học Bách Khoa TP.HCM']
      }
    ];

    // helper: map affiliation_id -> std_name để in ra dễ nhìn
    const affIdToName = {};
    affResult.rows.forEach((row) => {
      affIdToName[row.id] = row.std_name;
    });

    // Thực hiện chèn dữ liệu
    console.log('>> Seeding users (chi tiết từng tài khoản):');
    for (const user of users) {
      const query = `
        INSERT INTO users (username, full_name, password_hash, role, is_active, student_id, email, affiliation_id) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING user_id
      `;
      const values = [
        user.username, user.full_name, user.password_hash, user.role, 
        true, user.student_id, user.email, user.affiliation_id
      ];

      const inserted = await client.query(query, values);
      const insertedId = inserted?.rows?.[0]?.user_id;

      console.log(
        `   - user_id=${insertedId} | username=${user.username} | role=${user.role} | full_name=${user.full_name} | student_id=${user.student_id} | email=${user.email} | affiliation=${affIdToName[user.affiliation_id] || user.affiliation_id}`
      );
    }

    console.log(">> Đã seed dữ liệu thành công cho " + users.length + " người dùng.");

  } catch (err) {
    console.error("!! Lỗi khi thực hiện seed dữ liệu:", err.message);
  } finally {
    await client.end();
    console.log("--- Ngắt kết nối Database ---");
  }
}

seedData();