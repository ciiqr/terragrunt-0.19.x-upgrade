# run
```bash
go run . ~/some/terragrunt/file/terraform.tfvars ~/some/other/terragrunt/file/terraform.tfvars
```

# LIMITATIONS:
- this only does the first two step of the terragrunt 0.19 upgrade guide: https://github.com/gruntwork-io/terragrunt/blob/master/_docs/migration_guides/upgrading_to_terragrunt_0.19.x.md
- you will still have to follow that guide to:
	- update the hcl syntax
	- change any incompatible usages of attributes/blocks
	- rename changed functions
