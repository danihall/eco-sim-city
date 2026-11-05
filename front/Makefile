.PHONY: compile
compile:
	@tsc && node scripts/copy-css.js && node scripts/check-spritesheet-values.js  && node scripts/copy-spritesheet.js

.PHONY: compile-from-scratch
compile-from-scratch:
	@rm -rf dist && make compile
