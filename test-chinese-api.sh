#!/bin/bash

echo "========================================"
echo "🐠 测试中文参数 API 调用"
echo "========================================"
echo ""

BASE_URL="http://localhost:8080/api/v1"

echo "📝 测试说明："
echo "  以下示例展示了如何使用中文参数调用后端 API"
echo "  后端会自动将中文转换为对应的英文/数值"
echo ""

echo "========================================"
echo "1️⃣ 设备控制 - 中文参数示例"
echo "========================================"
echo ""
echo "路径参数支持中文："
echo "  /api/v1/device/灯光/主灯/toggle   → 自动转换为 /device/light/主灯/toggle"
echo "  /api/v1/device/气泵/主气泵/value  → 自动转换为 /device/pump/主气泵/value"
echo ""
echo "请求体支持中文："
cat << 'EOF'
  # 开关支持中文
  { "status": "开启" }    → true
  { "status": "关闭" }    → false
  { "status": "开" }      → true
  { "status": "关" }      → false
  
  # 模式支持中文
  { "manual_mode": "手动" }    → true
  { "manual_mode": "自动" }    → false
  { "manual_mode": "手动模式" } → true
  { "manual_mode": "自动模式" } → false
  
  # 灯光亮度支持中文
  { "current_value": "低光" } → 30
  { "current_value": "中光" } → 60
  { "current_value": "高光" } → 100
  { "current_value": "弱" }   → 30
  { "current_value": "强" }   → 100
  
  # 气泵档位支持中文
  { "current_value": "1档" } → 1
  { "current_value": "2档" } → 2
  { "current_value": "3档" } → 3
  { "current_value": "4档" } → 4
  { "current_value": "5档" } → 5
  { "current_value": "一档" } → 1
  { "current_value": "三档" } → 3
  { "current_value": "五档" } → 5
  { "current_value": "最低" } → 1
  { "current_value": "中等" } → 3
  { "current_value": "最高" } → 5
EOF
echo ""

echo "========================================"
echo "2️⃣ 光照排程 - 中文参数示例"
echo "========================================"
echo ""
cat << 'EOF'
  # 亮度支持中文
  { "brightness": "低" }  → 30
  { "brightness": "中" }  → 60
  { "brightness": "高" }  → 100
  
  # 启用状态支持中文
  { "enabled": "开启" } → true
  { "enabled": "关闭" } → false
  
  # 时间段支持中文别名
  { "time_range": "早上" } → start_time: "06:00:00", end_time: "09:00:00"
  { "time_range": "上午" } → start_time: "09:00:00", end_time: "12:00:00"
  { "time_range": "中午" } → start_time: "12:00:00", end_time: "14:00:00"
  { "time_range": "下午" } → start_time: "14:00:00", end_time: "18:00:00"
  { "time_range": "傍晚" } → start_time: "17:00:00", end_time: "19:00:00"
  { "time_range": "晚上" } → start_time: "19:00:00", end_time: "22:00:00"
  { "time_range": "夜间" } → start_time: "22:00:00", end_time: "06:00:00"
  { "time_range": "全天" } → start_time: "00:00:00", end_time: "23:59:59"
EOF
echo ""

echo "========================================"
echo "3️⃣ 溶氧量配置 - 中文参数示例"
echo "========================================"
echo ""
cat << 'EOF'
  # 气泵档位支持中文
  { "pump_level": "3档" } → 3
  { "pump_level": "五档" } → 5
  { "pump_level": "中等" } → 3
EOF
echo ""

echo "========================================"
echo "4️⃣ 实际 API 测试（需要后端服务启动）"
echo "========================================"
echo ""

# 检查后端是否启动
echo "🔍 检查后端服务是否启动..."
if curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/device/dashboard" | grep -q "200"; then
    echo "✅ 后端服务已启动，开始测试..."
    echo ""
    
    echo "📡 测试 1: 使用中文路径获取设备状态"
    echo "GET $BASE_URL/device/灯光/主灯"
    curl -s "$BASE_URL/device/灯光/主灯" | python3 -m json.tool
    echo ""
    
    echo "📡 测试 2: 使用中文状态切换设备"
    echo "PUT $BASE_URL/device/灯光/主灯/toggle"
    echo 'Body: {"status": "开启"}'
    curl -s -X PUT "$BASE_URL/device/灯光/主灯/toggle" \
      -H "Content-Type: application/json" \
      -d '{"status": "开启"}' | python3 -m json.tool
    echo ""
    
    echo "📡 测试 3: 使用中文亮度值"
    echo "PUT $BASE_URL/device/灯光/主灯/value"
    echo 'Body: {"current_value": "高光"}'
    curl -s -X PUT "$BASE_URL/device/灯光/主灯/value" \
      -H "Content-Type: application/json" \
      -d '{"current_value": "高光"}' | python3 -m json.tool
    echo ""
    
    echo "📡 测试 4: 使用中文档位值"
    echo "PUT $BASE_URL/device/气泵/主气泵/value"
    echo 'Body: {"current_value": "3档"}'
    curl -s -X PUT "$BASE_URL/device/气泵/主气泵/value" \
      -H "Content-Type: application/json" \
      -d '{"current_value": "3档"}' | python3 -m json.tool
    echo ""
    
    echo "📡 测试 5: 使用中文模式切换"
    echo "PUT $BASE_URL/device/灯光/主灯/manual"
    echo 'Body: {"manual_mode": "自动模式"}'
    curl -s -X PUT "$BASE_URL/device/灯光/主灯/manual" \
      -H "Content-Type: application/json" \
      -d '{"manual_mode": "自动模式"}' | python3 -m json.tool
    echo ""
    
    echo "📡 测试 6: 使用中文时间段创建光照排程"
    echo "POST $BASE_URL/lighting/schedules"
    echo 'Body: {"name": "测试排程", "time_range": "下午", "brightness": "中", "enabled": "开启"}'
    curl -s -X POST "$BASE_URL/lighting/schedules" \
      -H "Content-Type: application/json" \
      -d '{"name": "测试排程", "time_range": "下午", "brightness": "中", "enabled": "开启"}' | python3 -m json.tool
    echo ""
    
    echo "📡 测试 7: 使用中文参数进行溶氧量匹配"
    echo "POST $BASE_URL/oxygen/match"
    echo 'Body: {"light_wattage": "80", "temperature": "25.5"}'
    curl -s -X POST "$BASE_URL/oxygen/match" \
      -H "Content-Type: application/json" \
      -d '{"light_wattage": "80", "temperature": "25.5"}' | python3 -m json.tool
    echo ""
    
    echo "📡 测试 8: 使用中文档位创建溶氧量配置"
    echo "POST $BASE_URL/oxygen/configs"
    echo 'Body: {"min_light_wattage": 0, "max_light_wattage": 30, "min_temperature": 0, "max_temperature": 20, "pump_level": "2档", "description": "测试配置"}'
    curl -s -X POST "$BASE_URL/oxygen/configs" \
      -H "Content-Type: application/json" \
      -d '{"min_light_wattage": 0, "max_light_wattage": 30, "min_temperature": 0, "max_temperature": 20, "pump_level": "2档", "description": "测试配置"}' | python3 -m json.tool
    echo ""
    
else
    echo "⚠️  后端服务未启动，跳过实际测试"
    echo "   请先启动后端服务：cd backend && go run cmd/main.go"
fi
echo ""

echo "========================================"
echo "5️⃣ 支持的中文映射表"
echo "========================================"
echo ""
echo "设备类型映射："
echo "  灯光, 主灯, 辅灯  → light"
echo "  气泵, 主气泵, 氧气泵 → pump"
echo ""
echo "状态映射："
echo "  开启, 开, 打开, on, true → true"
echo "  关闭, 关, off, false    → false"
echo ""
echo "模式映射："
echo "  手动, 手动模式, manual → true"
echo "  自动, 自动模式, auto   → false"
echo ""
echo "亮度映射："
echo "  低, 弱, 低光 → 30%"
echo "  中, 中等, 中光 → 60%"
echo "  高, 强, 高光, 最高 → 100%"
echo ""
echo "气泵档位映射："
echo "  1档, 一档, 最低 → 1"
echo "  2档, 二档, 低 → 2"
echo "  3档, 三档, 中, 中等 → 3"
echo "  4档, 四档, 高 → 4"
echo "  5档, 五档, 最高 → 5"
echo ""
echo "时间段映射："
echo "  早上 → 06:00-09:00"
echo "  上午 → 09:00-12:00"
echo "  中午 → 12:00-14:00"
echo "  下午 → 14:00-18:00"
echo "  傍晚 → 17:00-19:00"
echo "  晚上 → 19:00-22:00"
echo "  夜间 → 22:00-06:00"
echo "  全天 → 00:00-23:59:59"
echo ""
echo "========================================"
echo "✅ 中文参数支持已配置完成！"
echo "========================================"
