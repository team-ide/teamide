CREATE TABLE `TM_VERSION` (
  `version` varchar(50) NOT NULL COMMENT '版本',
  `detail` text NOT NULL COMMENT '说明',
  `createTime` datetime NOT NULL COMMENT '创建时间',
  `updateTime` datetime DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='版本';