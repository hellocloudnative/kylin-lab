-- 开始初始化数据 ;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
SHOW VARIABLES LIKE 'time_zone';
SET time_zone = 'SYSTEM';
SET GLOBAL time_zone = '+00:00';

INSERT INTO `lab_virtualmachine` (
    `user_id`,
    `user_name`,
    `uuid`,
    `cpu_architecture`,
    `os_image`,
    `flavors`,
    `ip_address`,
    `network_name`,
    `duration`,
    `time_ofuse`,
    `apply_status`,
    `apply_time`,
    `status`,
    `vm_log`,
    `created_at`,
    `updated_at`
) VALUES (
             1,
             'admin',
             'bab254b7-5074-4d1d-bdb4-54335de251ac',
             'aarch64',
             'arm2303.qcow2',
             '2C-8g',
             '192.168.0.10',
             'vxlan',
             30,
             '2024-06-13 11:05:27-2024-06-13 11:35:27',
             0,
             '2024-06-13 11:05:27.165529+08:00',
             0,
             '系统自动消息: 审批通过',
             '2024-06-13 11:05:27.165529+08:00',
             '2024-06-13 11:05:27.165529+08:00'
         );


SET FOREIGN_KEY_CHECKS = 1;
-- 数据完成 ;