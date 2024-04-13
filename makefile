PROJECTS = ssh-compose

.PHONY: default
default:
	@# do nothing

.PHONY: build-all
build-all: $(foreach project,$(PROJECTS),bin/$(project))

bin/%: default
	@mkdir -p bin
	@$(MAKE) -C tools/$(@F) $@
	@cp tools/$(@F)/$@ $@

echo:
	@echo $(PROJECTS)