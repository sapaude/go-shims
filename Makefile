# 单版本策略 (Monorepo with Single Version):
# 优点： 简单，所有模块同步更新，易于理解和维护。
# 缺点： 即使只修改了一个小模块，也需要发布整个仓库的新版本，可能导致不必要的版本号跳跃。
VERSION_FILE = version

.PHONY: tag push

tag:
	@current_version=$$(cat $(VERSION_FILE)); \
	version_number=$${current_version#v}; \
	IFS=. read major minor patch <<< "$$version_number"; \
	new_patch=$$((patch + 1)); \
	new_version="v$${major}.$${minor}.$${new_patch}"; \
	echo "$$new_version" > $(VERSION_FILE); \
	git add .; \
	git commit -m "Bump tag version to $$new_version" > /dev/null; \
	git tag -a "$${new_version}" -m "--feat $${new_version}"; \
	echo "Created tag: $${new_version}"

push:
	# 推送tag
	@current_version=$$(cat $(VERSION_FILE)); \
	git push origin --tags

	# 推送 main 分支
	git push origin main

update: tag push