-- 开始初始化数据 ;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
SHOW VARIABLES LIKE 'time_zone';
SET time_zone = 'SYSTEM';
SET GLOBAL time_zone = '+00:00';

INSERT INTO `lab_virtualMachine` (
   `vm_id`, `user_id`, `cpu_architecture`, `os_type`, `os_image`, `machine_spec`,
    `ip_address`, `duration`, `status`, `vm_log`, `vnc_address`, `created_at`, `updated_at`
) VALUES
(1,1, 'x86_64', 'Desktop', 'Kylin-Desktop-V10-SP1-General-Release-2303-x86-64', '4C-8g ', '192.168.1.1','30', 0, 'Log data...', '192.168.1.1:5900', NOW(), NOW()),
(2,1, 'arm_64', 'Server', 'Kylin-Server-10-SP2-Release-Build09-20210524-arm64', '4C-8g', '192.168.1.10','30' , 1, 'Log data...', '192.168.1.10:5900', NOW(), NOW()),
(3,1, 'arm_64', 'Server', 'Kylin-Server-10-SP2-Release-Build09-20210524-arm64', '4C-8g', '192.168.1.11','30', 0, 'Log data...', '192.168.1.10:5900', NOW(), NOW()),
(4,1, 'arm_64', 'Server', 'Kylin-Server-10-SP2-Release-Build09-20210524-arm64', '4C-8g', '192.168.1.12','30', 1, 'Log data...', '192.168.1.10:5900', NOW(), NOW()),
(5,1, 'arm_64', 'Server', 'Kylin-Server-10-SP2-Release-Build09-20210524-arm64', '4C-8g', '192.168.1.13','30' , 0, 'Log data...', '192.168.1.10:5900', NOW(), NOW());



SET FOREIGN_KEY_CHECKS = 1;
-- 数据完成 ;