CREATE TABLE `TM_INSTALL` (
  `module` varchar(50) NOT NULL COMMENT '模块',
  `stage` varchar(50) NOT NULL COMMENT '阶段',
  `detail` text NOT NULL COMMENT '描述',
  `createTime` datetime NOT NULL COMMENT '创建时间',
  `updateTime` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`module`, `stage`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '安装';