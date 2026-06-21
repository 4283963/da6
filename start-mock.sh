#!/bin/bash

echo "========================================"
echo "🐠 水族箱智能控制系统 - 模拟数据脚本"
echo "========================================"
echo ""

BASE_URL="http://localhost:8080/api/v1"

echo "📊 正在生成模拟传感器数据..."
echo ""

# 生成 48 条模拟数据（每小时一条，共2天）
for i in {0..47}; do
  # 模拟时间：从48小时前到现在
  # 水温模拟：白天高，晚上低
  hour_of_day=$(( (i + $(date +%H) - 48 + 24 * 3) % 24 ))
  base_temp=22
  if [ $hour_of_day -ge 10 ] && [ $hour_of_day -le 16 ]; then
    # 白天中午温度高
    temp_offset=$(( RANDOM % 40 / 10 ))
    temp=$(echo "scale=1; 25 + $temp_offset" | bc)
  else
    # 晚上温度低
    temp_offset=$(( RANDOM % 30 / 10 ))
    temp=$(echo "scale=1; 22 + $temp_offset" | bc)
  fi

  # 灯光瓦数模拟：根据排程时间
  if [ $hour_of_day -ge 7 ] && [ $hour_of_day -lt 19 ]; then
    if [ $hour_of_day -ge 9 ] && [ $hour_of_day -lt 17 ]; then
      # 日间强光 80% = 80W
      light_wattage=80
    else
      # 早晚弱光 30-50%
      light_wattage=$(( 30 + RANDOM % 21 ))
    fi
  else
    # 夜间关灯
    light_wattage=0
  fi

  # 随机溶氧量
  do_value=$(echo "scale=2; 6 + $(( RANDOM % 20 )) / 10" | bc)

  echo "  [$i/47] 时间: $hour_of_day:00, 水温: ${temp}°C, 灯光: ${light_wattage}W, DO: ${do_value}mg/L"

  curl -s -X POST "$BASE_URL/sensor/data" \
    -H "Content-Type: application/json" \
    -d "{
      \"temperature\": $temp,
      \"light_wattage\": $light_wattage,
      \"dissolved_oxygen\": $do_value
    }" > /dev/null

  sleep 0.1
done

echo ""
echo "✅ 模拟数据生成完成！"
echo ""
echo "📈 当前系统状态："
echo ""
echo "--- 光照状态 ---"
curl -s "$BASE_URL/lighting/status" | python3 -m json.tool 2>/dev/null || echo "请确保后端服务已启动"
echo ""
echo "--- 设备状态 ---"
curl -s "$BASE_URL/device/dashboard" | python3 -m json.tool 2>/dev/null || echo "请确保后端服务已启动"
echo ""
echo "--- 传感器统计 ---"
curl -s "$BASE_URL/sensor/stats?hours=24" | python3 -m json.tool 2>/dev/null || echo "请确保后端服务已启动"
echo ""
echo "========================================"
echo "🎉 模拟数据已注入，可以刷新前端查看效果！"
echo "========================================"
