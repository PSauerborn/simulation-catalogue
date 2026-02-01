.PHONY: lint
lint:
	terraform fmt -recursive
	terraform validate
	terraform-docs markdown table . > TF_DOCS.md
