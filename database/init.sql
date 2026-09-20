CREATE TABLE IF NOT EXISTS users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(80) NOT NULL,
  major VARCHAR(120) NOT NULL,
  campus VARCHAR(40) NOT NULL,
  credit_score INT NOT NULL DEFAULT 80,
  credit_level VARCHAR(40) NOT NULL
);

CREATE TABLE IF NOT EXISTS skills (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  title VARCHAR(120) NOT NULL,
  category VARCHAR(40) NOT NULL,
  level_score INT NOT NULL,
  campus VARCHAR(40) NOT NULL,
  description TEXT NOT NULL,
  portfolio VARCHAR(160) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 技能未来一周可交换时隙：weekday(1=周一..7=周日) 与 part(1=上午,2=下午,3=晚上)
-- 同 weekday 内相邻 part 构成连续可交换时段。
CREATE TABLE IF NOT EXISTS skill_time_slots (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  skill_id BIGINT NOT NULL,
  weekday TINYINT NOT NULL,
  part TINYINT NOT NULL,
  label VARCHAR(20) NOT NULL,
  UNIQUE KEY uk_skill_slot (skill_id, weekday, part)
);

CREATE TABLE IF NOT EXISTS needs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  title VARCHAR(120) NOT NULL,
  category VARCHAR(40) NOT NULL,
  campus VARCHAR(40) NOT NULL,
  expect_time VARCHAR(80) NOT NULL,
  budget_type VARCHAR(40) NOT NULL,
  description TEXT NOT NULL
);

-- 需求方未来一周可用时隙，与 skill_time_slots 取交集判定共同时隙。
CREATE TABLE IF NOT EXISTS need_available_slots (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  need_id BIGINT NOT NULL,
  weekday TINYINT NOT NULL,
  part TINYINT NOT NULL,
  label VARCHAR(20) NOT NULL,
  UNIQUE KEY uk_need_slot (need_id, weekday, part)
);

-- 交换邀请：仅 pending 计入有效邀请（每个需求上限 3），accepted/expired 为终态。
CREATE TABLE IF NOT EXISTS invitations (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  need_id BIGINT NOT NULL,
  from_user VARCHAR(80) NOT NULL,
  to_user VARCHAR(80) NOT NULL,
  offer_skill VARCHAR(120) NOT NULL,
  wanted_skill VARCHAR(120) NOT NULL,
  campus VARCHAR(40) NOT NULL,
  proposed_slot VARCHAR(80) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  note VARCHAR(160) NOT NULL DEFAULT '',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  decided_at TIMESTAMP NULL,
  KEY idx_need_status (need_id, status)
);

CREATE TABLE IF NOT EXISTS appointments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  pair_name VARCHAR(120) NOT NULL,
  exchange_time VARCHAR(80) NOT NULL,
  place VARCHAR(120) NOT NULL,
  status VARCHAR(40) NOT NULL,
  agenda TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS reviews (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  from_user VARCHAR(80) NOT NULL,
  to_user VARCHAR(80) NOT NULL,
  rating INT NOT NULL,
  content TEXT NOT NULL
);

INSERT INTO users(name, major, campus, credit_score, credit_level) VALUES
('林澈', '新闻传播 2023', '东校区', 91, '黄金导师'),
('孟野', '音乐表演 2022', '东校区', 88, '白银协作者'),
('周芮', '统计学 2021', '中心校区', 93, '黄金导师'),
('许安', '社会学 2023', '中心校区', 85, '白银协作者'),
('韩沙', '视觉传达 2022', '西校区', 89, '白银协作者'),
('唐鹿', '法语 2021', '东校区', 76, '青铜互助者'),
('沈霜', '法学 2022', '东校区', 83, '白银协作者'),
('姜澜', '摄影 2020', '东校区', 87, '白银协作者');

INSERT INTO skills(user_id, title, category, level_score, campus, description, portfolio) VALUES
(1, '毕业照人像摄影', '摄影', 92, '东校区', '提供构图、修图和毕业季跟拍，可交换吉他入门课。', '12组校园人像作品'),
(2, '民谣吉他陪练', '乐器', 81, '东校区', '节奏型、弹唱和舞台经验分享，想找人拍宣传照。', '校园音乐节演出视频'),
(3, 'Python 数据分析', '编程', 88, '中心校区', 'pandas、可视化、论文数据清洗辅导。', '3份课程项目证书'),
(4, '乐队宣传照精修', '摄影', 84, '西校区', '擅长室外调色与海报精修，可跨校区接单。', '2支乐队海报作品'),
(6, '法语口语陪练', '外语', 79, '东校区', 'DELF B2 口语训练和发音纠正。', 'DELF B2 证书'),
(7, '法语论文润色', '外语', 90, '东校区', '法学法语双语，擅长学术写作与答辩口语。', '2篇法语期刊润色记录'),
(8, '活动纪实摄影', '摄影', 86, '东校区', '社团活动与宣传照纪实拍摄。', '8场活动纪实相册');

-- 林澈：周三晚、周六上午、周六下午（周六上午+下午为连续时隙）
INSERT INTO skill_time_slots(skill_id, weekday, part, label) VALUES
(1, 3, 3, '周三晚上'), (1, 6, 1, '周六上午'), (1, 6, 2, '周六下午');
-- 孟野：周三晚、周六上午、周六下午
INSERT INTO skill_time_slots(skill_id, weekday, part, label) VALUES
(2, 3, 3, '周三晚上'), (2, 6, 1, '周六上午'), (2, 6, 2, '周六下午');

INSERT INTO needs(user_id, title, category, campus, expect_time, budget_type, description) VALUES
(2, '找人帮忙拍乐队宣传照', '摄影', '东校区', '本周六下午', '技能交换', '可交换 3 次吉他课，希望会调色和室外构图。'),
(4, '求教 Python 数据分析', '编程', '中心校区', '周二晚上', '小额报酬', '论文问卷数据需要清洗和画图，最好有 pandas 经验。'),
(1, '想学吉他扫弦入门', '乐器', '东校区', '周三晚上', '技能交换', '用摄影课交换吉他基础，希望同校区或线上。'),
(1, '求法语发音陪练', '外语', '东校区', '周六下午', '技能交换', '准备 DELF 入门，希望同校区同学帮忙纠音。');

-- 需求 1（孟野）可用时隙：周三晚、周六上午、周六下午
INSERT INTO need_available_slots(need_id, weekday, part, label) VALUES
(1, 3, 3, '周三晚上'), (1, 6, 1, '周六上午'), (1, 6, 2, '周六下午');

-- 预置两条待确认邀请：任一被接受后另一条应立即失效。
INSERT INTO invitations(need_id, from_user, to_user, offer_skill, wanted_skill, campus, proposed_slot, status) VALUES
(1, '林澈', '孟野', '毕业照人像摄影', '民谣吉他陪练', '东校区', '周六下午', 'pending'),
(1, '姜澜', '孟野', '活动纪实摄影', '民谣吉他陪练', '东校区', '周六下午', 'pending');
