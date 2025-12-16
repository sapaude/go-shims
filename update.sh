#!/bin/bash

set -e

# 获取仓库根目录
REPO_ROOT=$(git rev-parse --show-toplevel)
VERSION_FILE="$REPO_ROOT/version"

# 读取当前版本
CURRENT_VERSION=$(cat "$VERSION_FILE" 2>/dev/null | tr -d '[:space:]' || echo "v1.0.0")

# 解析版本号 (v1.2.3 -> 1 2 3)
parse_version() {
    local ver=${1#v}
    echo $ver | tr '.' ' '
}

# 版本迭代
bump_version() {
    local type=${1:-patch}
    read -r major minor patch <<< $(parse_version "$CURRENT_VERSION")

    case $type in
        major) echo "v$((major + 1)).0.0" ;;
        minor) echo "v${major}.$((minor + 1)).0" ;;
        patch) echo "v${major}.${minor}.$((patch + 1))" ;;
        *) echo "Error: Invalid type '$type'. Use: major|minor|patch" >&2; exit 1 ;;
    esac
}

# 使用说明
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    cat << EOF
Usage: $0 [major|minor|patch|v1.0.5]

  major         $CURRENT_VERSION -> $(bump_version major)
  minor         $CURRENT_VERSION -> $(bump_version minor)
  patch         $CURRENT_VERSION -> $(bump_version patch) (default)
  v1.0.5        指定版本号

Current: $CURRENT_VERSION
EOF
    exit 0
fi

# 计算新版本
if [[ $1 =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    NEW_VERSION=$1
else
    NEW_VERSION=$(bump_version ${1:-patch})
fi

echo "Version: $CURRENT_VERSION -> $NEW_VERSION"

# 执行更新
echo "$NEW_VERSION" > "$VERSION_FILE"
git add .
git commit -m "chore: bump version to $NEW_VERSION [x/log]"
git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"

# 推送
BRANCH=$(git rev-parse --abbrev-ref HEAD)
git push origin "$BRANCH"
git push origin "$NEW_VERSION"

echo "✓ Updated to $NEW_VERSION"
