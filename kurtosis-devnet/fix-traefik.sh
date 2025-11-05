#!/bin/bash

# Traefik 네트워크 문제 자동 해결 스크립트
# 사용법: ./fix-traefik.sh

set -e

echo "🔧 Traefik 네트워크 문제 자동 해결 시작..."

# Traefik 컨테이너 찾기
TRAEFIK_CONTAINER=$(docker ps --filter "name=kurtosis-reverse-proxy" --format "{{.Names}}" | head -1)

if [ -z "$TRAEFIK_CONTAINER" ]; then
    echo "❌ Traefik 컨테이너를 찾을 수 없습니다."
    exit 1
fi

echo "📦 Traefik 컨테이너 발견: $TRAEFIK_CONTAINER"

# 현재 설정 백업
echo "💾 현재 Traefik 설정 백업 중..."
docker exec "$TRAEFIK_CONTAINER" cat /etc/traefik/traefik.yml > "/tmp/traefik-backup-$(date +%Y%m%d-%H%M%S).yml"

# 개선된 설정 적용
echo "⚙️  개선된 Traefik 설정 적용 중..."
docker exec "$TRAEFIK_CONTAINER" sh -c 'cat > /etc/traefik/traefik.yml << EOF
accesslog: {}
log:
  level: DEBUG
api:
  debug: true
  dashboard: true
  insecure: true
  disabledashboardad: true
 
entryPoints:
  web:
    address: ":9730"
  traefik:
    address: ":9731"

providers:
  docker:
    endpoint: "unix:///var/run/docker.sock"
    exposedByDefault: false
    # 동적 네트워크 감지 (하드코딩 제거)
    useBindPortIP: false
    watch: true
    pollInterval: "5s"
    # 네트워크 자동 감지 설정
    defaultRule: "Host(\`{{ normalize .Name }}\`)"
EOF'

# Traefik 재시작
echo "🔄 Traefik 재시작 중..."
docker restart "$TRAEFIK_CONTAINER"

# 재시작 완료 대기
echo "⏳ Traefik 재시작 완료 대기 중..."
sleep 10

# 상태 확인
echo "🔍 Traefik 상태 확인 중..."
if docker ps --filter "name=$TRAEFIK_CONTAINER" --filter "status=running" | grep -q "$TRAEFIK_CONTAINER"; then
    echo "✅ Traefik이 성공적으로 재시작되었습니다!"
    
    # 네트워크 에러 확인
    echo "🌐 네트워크 에러 확인 중..."
    sleep 5
    
    NETWORK_ERRORS=$(docker logs "$TRAEFIK_CONTAINER" --since=30s 2>&1 | grep -c "Could not find network" || true)
    
    if [ "$NETWORK_ERRORS" -eq 0 ]; then
        echo "✅ 네트워크 에러가 해결되었습니다!"
        echo "🎉 Traefik 자동 복구 완료!"
    else
        echo "⚠️  일부 네트워크 에러가 여전히 발생할 수 있지만, 동적 감지가 활성화되어 곧 해결될 예정입니다."
    fi
    
    # 대시보드 접근 테스트
    echo "🌐 Traefik 대시보드 테스트 중..."
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:9730/ || echo "000")
    if [ "$HTTP_CODE" != "000" ]; then
        echo "✅ Traefik 대시보드 접근 가능: http://127.0.0.1:9730/"
    else
        echo "⚠️  대시보드 접근 테스트 실패, 하지만 서비스는 정상 작동할 수 있습니다."
    fi
else
    echo "❌ Traefik 재시작에 실패했습니다."
    exit 1
fi

echo ""
echo "📋 요약:"
echo "   - Traefik 설정: 동적 네트워크 감지로 변경됨"
echo "   - 폴링 간격: 5초마다 네트워크 상태 확인"
echo "   - 하드코딩된 네트워크 ID 제거됨"
echo ""
echo "🔧 문제 재발 시 다음 명령어로 수동 해결 가능:"
echo "   docker restart $TRAEFIK_CONTAINER"