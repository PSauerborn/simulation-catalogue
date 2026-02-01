.PHONY: scan-secrets
scan-secrets:
	detect-secrets scan \
		--exclude-files .*/tests/fixtures \
		. > .secrets.baseline
	detect-secrets audit .secrets.baseline
