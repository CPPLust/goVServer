


-- 1. 创建数据库（如不存在）
CREATE DATABASE IF NOT EXISTS video_server
    DEFAULT CHARACTER SET utf8mb4;


-- 2. 切换到目标数据库
USE video_server;


-- 3. 创建 users 表（仅当表不存在时创建）
CREATE TABLE IF NOT EXISTS users (
     id         INT           NOT NULL AUTO_INCREMENT COMMENT '用户ID',
     login_name VARCHAR(64)   NOT NULL                COMMENT '登录名',
     pwd        TEXT          NOT NULL                COMMENT '密码（建议存储哈希值）',
	 PRIMARY KEY (id),
	 UNIQUE KEY `uk_login_name` (`login_name`)
) DEFAULT CHARSET=utf8mb4
  COMMENT='用户表';
  
  
CREATE TABLE IF NOT EXISTS video_info (
  id VARCHAR(255) PRIMARY KEY NOT NULL,
  author_id INT,
  name TEXT,
  display_ctime TEXT,
  create_time DATETIME
);

CREATE TABLE IF NOT EXISTS comments (
  id VARCHAR(64) PRIMARY KEY NOT NULL,
  video_id VARCHAR(64),
  author_id INT,
  content TEXT,
  time DATETIME
);

CREATE TABLE IF NOT EXISTS sessions (
  session_id VARCHAR(255) PRIMARY KEY NOT NULL,
  TTL TINYTEXT,
  login_name VARCHAR(255)
);