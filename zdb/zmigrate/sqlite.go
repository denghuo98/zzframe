package zmigrate

// sqlite 建表语句

var createAdminMemberTableSQLite = `
-- 系统管理员表
CREATE TABLE IF NOT EXISTS zz_admin_member (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  dept_id INTEGER DEFAULT 0,
  real_name TEXT DEFAULT '',
  username TEXT NOT NULL DEFAULT '',
  password_hash TEXT NOT NULL DEFAULT '',
  salt TEXT NOT NULL,
  password_reset_token TEXT DEFAULT '',
  avatar TEXT DEFAULT '',
  sex INTEGER DEFAULT 3,
  email TEXT DEFAULT '',
  mobile TEXT DEFAULT '',
  last_active_at TEXT,
  remark TEXT,
  status INTEGER DEFAULT 1,
  created_at TEXT,
  updated_at TEXT
);
`

var createAdminMemberRoleTableSQLite = `
-- 管理员角色关联表
CREATE TABLE IF NOT EXISTS zz_admin_member_role (
  member_id INTEGER NOT NULL,
  role_id INTEGER NOT NULL,
  PRIMARY KEY (member_id, role_id)
);
`

var createAdminMenuTableSQLite = `
CREATE TABLE IF NOT EXISTS zz_admin_menu (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  pid INTEGER DEFAULT 0,
  title TEXT NOT NULL,
  name TEXT NOT NULL,
  path TEXT,
  icon TEXT,
  type INTEGER NOT NULL DEFAULT 1,
  redirect TEXT,
  component TEXT NOT NULL,
  is_frame INTEGER DEFAULT 1,
  frame_src TEXT,
  hidden INTEGER DEFAULT 0,
  sort INTEGER DEFAULT 0,
  remark TEXT,
  status INTEGER DEFAULT 1,
  updated_at TEXT,
  created_at TEXT,
  permissions TEXT NOT NULL DEFAULT ''
);
`

var createAdminRoleTableSQLite = `
CREATE TABLE IF NOT EXISTS zz_admin_role (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  key TEXT NOT NULL,
  remark TEXT,
  sort INTEGER NOT NULL DEFAULT 0,
  status INTEGER NOT NULL DEFAULT 1,
  created_at TEXT,
  updated_at TEXT
);
`

var createAdminRoleCasbinTableSQLite = `
CREATE TABLE IF NOT EXISTS zz_admin_role_casbin (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  p_type TEXT,
  v0 TEXT,
  v1 TEXT,
  v2 TEXT,
  v3 TEXT,
  v4 TEXT,
  v5 TEXT
);
`

var createAdminRoleMenuTableSQLite = `
CREATE TABLE IF NOT EXISTS zz_admin_role_menu (
  role_id INTEGER NOT NULL,
  menu_id INTEGER NOT NULL,
  PRIMARY KEY (role_id, menu_id)
);
`

var createSysLoginLogTableSQLite = `
CREATE TABLE IF NOT EXISTS zz_sys_login_log (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  req_id TEXT,
  member_id INTEGER DEFAULT 0,
  username TEXT,
  response TEXT,
  login_at TEXT,
  login_ip TEXT,
  province TEXT,
  city TEXT,
  user_agent TEXT,
  err_msg TEXT,
  status INTEGER NOT NULL DEFAULT 1,
  created_at TEXT,
  updated_at TEXT
);
`

var createSysAttachmentTableSQLite = `
CREATE TABLE IF NOT EXISTS zz_sys_attachment (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  app_id TEXT NOT NULL,
  member_id INTEGER DEFAULT 0,
  cate_id INTEGER DEFAULT 0,
  drive TEXT,
  name TEXT,
  kind TEXT,
  mime_type TEXT NOT NULL DEFAULT '',
  naive_type TEXT,
  path TEXT,
  file_url TEXT,
  size INTEGER DEFAULT 0,
  ext TEXT,
  md5 TEXT,
  status INTEGER NOT NULL DEFAULT 1,
  created_at TEXT DEFAULT (datetime('now','localtime')),
  updated_at TEXT DEFAULT (datetime('now','localtime'))
);
`
