-- 管理员_用户表
CREATE TABLE `zz_admin_member` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `dept_id` INTEGER DEFAULT 0,
  `real_name` TEXT DEFAULT '',
  `username` TEXT NOT NULL DEFAULT '',
  `password_hash` TEXT NOT NULL DEFAULT '',
  `salt` TEXT NOT NULL,
  `password_reset_token` TEXT DEFAULT '',
  `avatar` TEXT DEFAULT '',
  `sex` INTEGER DEFAULT 3,
  `email` TEXT DEFAULT '',
  `mobile` TEXT DEFAULT '',
  `last_active_at` TEXT DEFAULT NULL,
  `remark` TEXT DEFAULT NULL,
  `status` INTEGER DEFAULT 1,
  `created_at` TEXT DEFAULT NULL,
  `updated_at` TEXT DEFAULT NULL
);

CREATE INDEX `idx_username` ON `zz_admin_member` (`username`);
CREATE INDEX `idx_phone` ON `zz_admin_member` (`mobile`);
CREATE INDEX `idx_email` ON `zz_admin_member` (`email`);


-- 管理员_用户角色关联
CREATE TABLE `zz_admin_member_role` (
  `member_id` INTEGER NOT NULL,
  `role_id` INTEGER NOT NULL,
  PRIMARY KEY (`member_id`,`role_id`)
);

CREATE INDEX `idx_member_id` ON `zz_admin_member_role` (`member_id`);
CREATE INDEX `idx_role_id` ON `zz_admin_member_role` (`role_id`);


-- 管理员_菜单权限
CREATE TABLE `zz_admin_menu` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `pid` INTEGER DEFAULT 0,
  `title` TEXT NOT NULL,
  `name` TEXT NOT NULL,
  `path` TEXT DEFAULT NULL,
  `icon` TEXT DEFAULT NULL,
  `type` INTEGER NOT NULL DEFAULT 1,
  `redirect` TEXT DEFAULT NULL,
  `component` TEXT NOT NULL,
  `is_frame` INTEGER DEFAULT 1,
  `frame_src` TEXT DEFAULT NULL,
  `hidden` INTEGER DEFAULT 0,
  `sort` INTEGER DEFAULT 0,
  `remark` TEXT DEFAULT NULL,
  `status` INTEGER DEFAULT 1,
  `updated_at` TEXT DEFAULT NULL,
  `created_at` TEXT DEFAULT NULL,
  `permissions` TEXT NOT NULL DEFAULT ''
);


-- 管理员_角色信息
CREATE TABLE `zz_admin_role` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `name` TEXT NOT NULL,
  `key` TEXT NOT NULL,
  `remark` TEXT DEFAULT NULL,
  `sort` INTEGER NOT NULL DEFAULT 0,
  `status` INTEGER NOT NULL DEFAULT 1,
  `created_at` TEXT DEFAULT NULL,
  `updated_at` TEXT DEFAULT NULL
);

CREATE INDEX `idx_name` ON `zz_admin_role` (`name`);
CREATE INDEX `idx_key` ON `zz_admin_role` (`key`);


-- 管理员_casbin权限表
CREATE TABLE `zz_admin_role_casbin` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `p_type` TEXT DEFAULT NULL,
  `v0` TEXT DEFAULT NULL,
  `v1` TEXT DEFAULT NULL,
  `v2` TEXT DEFAULT NULL,
  `v3` TEXT DEFAULT NULL,
  `v4` TEXT DEFAULT NULL,
  `v5` TEXT DEFAULT NULL
);


-- 管理员_角色菜单关联
CREATE TABLE `zz_admin_role_menu` (
  `role_id` INTEGER NOT NULL,
  `menu_id` INTEGER NOT NULL,
  PRIMARY KEY (`role_id`,`menu_id`)
);

CREATE INDEX `idx_role_id` ON `zz_admin_role_menu` (`role_id`);
CREATE INDEX `idx_menu_id` ON `zz_admin_role_menu` (`menu_id`);


-- 系统_登录日志
CREATE TABLE `zz_sys_login_log` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `req_id` TEXT DEFAULT NULL,
  `member_id` INTEGER DEFAULT 0,
  `username` TEXT DEFAULT NULL,
  `response` TEXT DEFAULT NULL,
  `login_at` TEXT DEFAULT NULL,
  `login_ip` TEXT DEFAULT NULL,
  `province` TEXT DEFAULT NULL,
  `city` TEXT DEFAULT NULL,
  `user_agent` TEXT DEFAULT NULL,
  `err_msg` TEXT DEFAULT NULL,
  `status` INTEGER NOT NULL DEFAULT 1,
  `created_at` TEXT DEFAULT NULL,
  `updated_at` TEXT DEFAULT NULL
);


-- 系统_附件管理
CREATE TABLE `zz_sys_attachment` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `app_id` TEXT NOT NULL,
  `member_id` INTEGER DEFAULT 0,
  `cate_id` INTEGER DEFAULT 0,
  `drive` TEXT DEFAULT NULL,
  `name` TEXT DEFAULT NULL,
  `kind` TEXT DEFAULT NULL,
  `mime_type` TEXT NOT NULL DEFAULT '',
  `naive_type` TEXT DEFAULT NULL,
  `path` TEXT DEFAULT NULL,
  `file_url` TEXT DEFAULT NULL,
  `size` INTEGER DEFAULT 0,
  `ext` TEXT DEFAULT NULL,
  `md5` TEXT DEFAULT NULL,
  `status` INTEGER NOT NULL DEFAULT 1,
  `created_at` TEXT DEFAULT (datetime('now','localtime')),
  `updated_at` TEXT DEFAULT (datetime('now','localtime'))
);
