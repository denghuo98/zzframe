package zmigrate

// mysql 建表语句

var createAdminMemberTableSQL = `
-- 系统管理员表
CREATE TABLE zz_admin_member (
  id bigint NOT NULL AUTO_INCREMENT COMMENT '管理员ID',
  dept_id bigint DEFAULT '0' COMMENT '部门ID',
  real_name varchar(32) DEFAULT '' COMMENT '真实姓名',
  username varchar(20) NOT NULL DEFAULT '' COMMENT '帐号',
  password_hash char(32) NOT NULL DEFAULT '' COMMENT '密码',
  salt char(16) NOT NULL COMMENT '密码盐',
  password_reset_token varchar(150) DEFAULT '' COMMENT '密码重置令牌',
  avatar char(150) DEFAULT '' COMMENT '头像',
  sex tinyint(1) DEFAULT '3' COMMENT '性别',
  email varchar(60) DEFAULT '' COMMENT '邮箱',
  mobile varchar(20) DEFAULT '' COMMENT '手机号码',
  last_active_at datetime DEFAULT NULL COMMENT '最后活跃时间',
  remark varchar(255) DEFAULT NULL COMMENT '备注',
  status tinyint(1) DEFAULT '1' COMMENT '状态',
  created_at datetime DEFAULT NULL COMMENT '创建时间',
  updated_at datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (id) USING BTREE,
  KEY idx_username (username),
  KEY idx_phone (mobile),
  KEY idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理员_用户表';
`

var createAdminMemberRoleTableSQL = `
-- 管理员角色关联表
CREATE TABLE zz_admin_member_role (
  member_id bigint NOT NULL COMMENT '管理员ID',
  role_id bigint NOT NULL COMMENT '角色ID',
  PRIMARY KEY (member_id, role_id) USING BTREE,
  KEY idx_member_id (member_id),
  KEY idx_role_id (role_id)
  
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理员_用户角色关联';
`

var createAdminMenuTableSql = `
CREATE TABLE zz_admin_menu (
  id bigint NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  pid bigint DEFAULT '0' COMMENT '父菜单ID',
  title varchar(64) NOT NULL COMMENT '菜单名称',
  name varchar(128) NOT NULL COMMENT '名称编码',
  path varchar(200) DEFAULT NULL COMMENT '路由地址',
  icon varchar(128) DEFAULT NULL COMMENT '菜单图标',
  type tinyint(1) NOT NULL DEFAULT '1' COMMENT '菜单类型（1目录 2菜单 3按钮）',
  redirect varchar(255) DEFAULT NULL COMMENT '重定向地址',
  component varchar(255) NOT NULL COMMENT '组件路径',
  is_frame tinyint(1) DEFAULT '1' COMMENT '是否内嵌',
  frame_src varchar(512) DEFAULT NULL COMMENT '内联外部地址',
  hidden tinyint(1) DEFAULT '0' COMMENT '是否隐藏',
  sort int DEFAULT '0' COMMENT '排序',
  remark varchar(255) DEFAULT NULL COMMENT '备注',
  status tinyint(1) DEFAULT '1' COMMENT '菜单状态',
  updated_at datetime DEFAULT NULL COMMENT '更新时间',
  created_at datetime DEFAULT NULL COMMENT '创建时间',
  permissions varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '菜单权限',
  PRIMARY KEY (id) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=129 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理员_菜单权限';
`

var createAdminRoleTableSQL = `
CREATE TABLE zz_admin_role (
  id bigint NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  name varchar(32) NOT NULL COMMENT '角色名称',
  key varchar(128) NOT NULL COMMENT '角色权限字符串',
  remark varchar(255) DEFAULT NULL COMMENT '备注',
  sort int NOT NULL DEFAULT '0' COMMENT '排序',
  status tinyint(1) NOT NULL DEFAULT '1' COMMENT '角色状态',
  created_at datetime DEFAULT NULL COMMENT '创建时间',
  updated_at datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (id) USING BTREE,
  KEY idx_name (name),
  KEY idx_key (key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理员_角色信息';
`

var createAdminRoleCasbinTableSQL = `
CREATE TABLE zz_admin_role_casbin (
  id bigint NOT NULL AUTO_INCREMENT COMMENT 'ID',
  p_type varchar(64) DEFAULT NULL,
  v0 varchar(256) DEFAULT NULL,
  v1 varchar(256) DEFAULT NULL,
  v2 varchar(256) DEFAULT NULL,
  v3 varchar(256) DEFAULT NULL,
  v4 varchar(256) DEFAULT NULL,
  v5 varchar(256) DEFAULT NULL,
  PRIMARY KEY (id) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='管理员_casbin权限表';
`

var createAdminRoleMenuTableSQL = `
CREATE TABLE zz_admin_role_menu (
  role_id bigint NOT NULL COMMENT '角色ID',
  menu_id bigint NOT NULL COMMENT '菜单ID',
  PRIMARY KEY (role_id, menu_id) USING BTREE,
  KEY idx_role_id (role_id),
  KEY idx_menu_id (menu_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理员_角色菜单关联';
`

var createSysLoginLogTableSQL = `
CREATE TABLE zz_sys_login_log (
  id bigint NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  req_id varchar(50) DEFAULT NULL COMMENT '请求ID',
  member_id bigint DEFAULT '0' COMMENT '用户ID',
  username varchar(64) DEFAULT NULL COMMENT '用户名',
  response json DEFAULT NULL COMMENT '响应数据',
  login_at datetime DEFAULT NULL COMMENT '登录时间',
  login_ip varchar(128) DEFAULT NULL COMMENT '登录IP',
  province varchar(128) DEFAULT NULL COMMENT 'IP定位省份',
  city varchar(128) DEFAULT NULL COMMENT 'IP定位城市',
  user_agent varchar(512) DEFAULT NULL COMMENT 'UA信息',
  err_msg varchar(1000) DEFAULT NULL COMMENT '错误提示',
  status tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  created_at datetime DEFAULT NULL COMMENT '创建时间',
  updated_at datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (id) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统_登录日志';
`

var createSysAttachmentTableSQL = `
CREATE TABLE zz_sys_attachment (
  id bigint NOT NULL AUTO_INCREMENT COMMENT '文件ID',
  app_id varchar(64) NOT NULL COMMENT '应用ID',
  member_id bigint DEFAULT '0' COMMENT '管理员ID',
  cate_id bigint unsigned DEFAULT '0' COMMENT '上传分类',
  drive varchar(64) DEFAULT NULL COMMENT '上传驱动',
  name varchar(1000) DEFAULT NULL COMMENT '文件原始名',
  kind varchar(16) DEFAULT NULL COMMENT '上传类型',
  mime_type varchar(128) NOT NULL DEFAULT '' COMMENT '扩展类型',
  naive_type varchar(32) DEFAULT NULL COMMENT 'NaiveUI类型',
  path varchar(1000) DEFAULT NULL COMMENT '本地路径',
  file_url varchar(1000) DEFAULT NULL COMMENT 'url',
  size bigint DEFAULT '0' COMMENT '文件大小',
  ext varchar(50) DEFAULT NULL COMMENT '扩展名',
  md5 varchar(32) DEFAULT NULL COMMENT 'md5校验码',
  status tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  created_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (id) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统_附件管理';
`
