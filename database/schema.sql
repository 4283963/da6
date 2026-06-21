CREATE DATABASE IF NOT EXISTS aquarium_control DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE aquarium_control;

-- 光照排程表
CREATE TABLE IF NOT EXISTS light_schedules (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '排程名称',
    start_time TIME NOT NULL COMMENT '开灯时间',
    end_time TIME NOT NULL COMMENT '关灯时间',
    brightness INT NOT NULL COMMENT '亮度百分比 0-100',
    enabled TINYINT(1) DEFAULT 1 COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_enabled (enabled),
    INDEX idx_time (start_time, end_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='光照排程表';

-- 溶氧量配置表
CREATE TABLE IF NOT EXISTS oxygen_configs (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    min_light_wattage INT NOT NULL COMMENT '最低灯光瓦数',
    max_light_wattage INT NOT NULL COMMENT '最高灯光瓦数',
    min_temperature DECIMAL(4,1) NOT NULL COMMENT '最低水温',
    max_temperature DECIMAL(4,1) NOT NULL COMMENT '最高水温',
    pump_level INT NOT NULL COMMENT '气泵档位 1-5',
    description VARCHAR(255) COMMENT '配置说明',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_light_range (min_light_wattage, max_light_wattage),
    INDEX idx_temp_range (min_temperature, max_temperature)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='溶氧量匹配配置表';

-- 设备状态表
CREATE TABLE IF NOT EXISTS device_status (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    device_type VARCHAR(50) NOT NULL COMMENT '设备类型: light/pump',
    device_name VARCHAR(100) NOT NULL COMMENT '设备名称',
    status TINYINT(1) DEFAULT 0 COMMENT '状态: 0关 1开',
    current_value INT COMMENT '当前值: 亮度/档位',
    manual_mode TINYINT(1) DEFAULT 0 COMMENT '手动模式: 0自动 1手动',
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_device (device_type, device_name),
    INDEX idx_device_type (device_type),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='设备状态表';

-- 环境传感器数据表
CREATE TABLE IF NOT EXISTS sensor_data (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    temperature DECIMAL(4,1) NOT NULL COMMENT '水温',
    light_wattage INT NOT NULL COMMENT '当前灯光瓦数',
    dissolved_oxygen DECIMAL(4,2) COMMENT '实测溶氧量',
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_recorded_at (recorded_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='环境传感器数据表';

INSERT INTO device_status (device_type, device_name, status, current_value, manual_mode) VALUES
('light', '主灯', 0, 0, 0),
('pump', '主气泵', 0, 1, 0),
('light', '辅灯', 0, 0, 0);

INSERT INTO oxygen_configs (min_light_wattage, max_light_wattage, min_temperature, max_temperature, pump_level, description) VALUES
(0, 30, 18.0, 24.0, 1, '低光低温'),
(0, 30, 24.1, 28.0, 2, '低光中温'),
(0, 30, 28.1, 32.0, 3, '低光高温'),
(31, 60, 18.0, 24.0, 2, '中光低温'),
(31, 60, 24.1, 28.0, 3, '中光中温'),
(31, 60, 28.1, 32.0, 4, '中光高温'),
(61, 100, 18.0, 24.0, 3, '高光低温'),
(61, 100, 24.1, 28.0, 4, '高光中温'),
(61, 100, 28.1, 32.0, 5, '高光高温');

INSERT INTO light_schedules (name, start_time, end_time, brightness, enabled) VALUES
('早间弱光', '07:00:00', '09:00:00', 30, 1),
('日间强光', '09:00:00', '17:00:00', 80, 1),
('傍晚渐暗', '17:00:00', '19:00:00', 50, 1);
