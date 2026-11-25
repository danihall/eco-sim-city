.PHONY: compile
compile:
	@tsc && \
	node scripts/copy-css.js && \
	node scripts/copy-fonts.js && \
	node scripts/check-spritesheet-values.js && \
	node scripts/copy-spritesheet.js

.PHONY: compile-from-scratch
compile-from-scratch:
	@rm -rf dist && make compile




.PHONY: compile-udon
compile-udon:
	esbuild file1.js file2.js file3.js \
	  --bundle \
	  --minify \
	  --outfile=bundle.min.js
