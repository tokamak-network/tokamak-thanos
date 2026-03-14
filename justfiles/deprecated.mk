# Deprecated Makefile targets - migrate to justfiles instead
.PHONY: $(DEPRECATED_TARGETS)
$(DEPRECATED_TARGETS):
	@echo "Target '$@' is deprecated. Please use justfiles instead."
	@echo "Run: just help"
	@exit 1
