
-- create database opcenter_config; 创建运营中心配置库

-- 用户节点信息表
CREATE TABLE IF NOT EXISTS `t_node` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `name` char(36) NOT NULL COMMENT '节点名称',
    `path` char(64) NOT NULL COMMENT '域名路径',
    `api_key` char(64) NOT NULL DEFAULT '' COMMENT '访问key',
    `api_secret` char(64) NOT NULL DEFAULT '' COMMENT '访问密钥',
    `created_on` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
    `created_by` char(36) DEFAULT NULL COMMENT '记录创建者ID',
    `modified_on` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录修改时间',
    `modified_by` char(36) DEFAULT NULL COMMENT '记录修改者ID',
    PRIMARY KEY (`id`),
    UNIQUE KEY `u_idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户节点管理表';

-- 用户统计指标表
CREATE TABLE IF NOT EXISTS `t_indicators` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `module` char(64) NOT NULL COMMENT '指标模块',
    `object` char(64) NOT NULL COMMENT '指标对象',
    `created_on` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
    `created_by` char(36) DEFAULT NULL COMMENT '记录创建者ID',
    `modified_on` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录修改时间',
    `modified_by` char(36) DEFAULT NULL COMMENT '记录修改者ID',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='统计指标项管理表';

-- 定时任务执行表
CREATE TABLE IF NOT EXISTS `t_crontab` (
    `id` char(36) NOT NULL COMMENT '主键ID',
    `name` varchar(64) NOT NULL COMMENT '定时任务名称',
    `schedule` varchar(64) NOT NULL COMMENT '定时任务执行策略',
    `expire` varchar(64) NOT NULL COMMENT '设定任务超时时间',
    `is_enable` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用,0-禁用,1-启用',
    `created_on` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
    `created_by` char(36) DEFAULT NULL COMMENT '记录创建者ID',
    `modified_on` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录修改时间',
    `modified_by` char(36) DEFAULT NULL COMMENT '记录修改者ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `u_idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='定时任务配置表';
