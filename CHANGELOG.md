# Changelog

<a name="unreleased"></a>
## [Unreleased]

- Upgrade for Terraform v0.12.0 ([#78](https://github.com/cloud-ca/terraform-provider-cloudca/issues/78))
- Use golangci lint instead of gometalinter ([#77](https://github.com/cloud-ca/terraform-provider-cloudca/issues/77))
- Compile provider with static binaries ([#76](https://github.com/cloud-ca/terraform-provider-cloudca/issues/76))


<a name="v1.4.1"></a>
## [v1.4.1] - 2019-04-17

- Fixed a bug which results in errors reading PF Rules, LB Rules and Network ACLs ([#75](https://github.com/cloud-ca/terraform-provider-cloudca/issues/75))


<a name="v1.4.0"></a>
## [v1.4.0] - 2019-04-15
### Documentation
- indicate input arguments are also returned ([#73](https://github.com/cloud-ca/terraform-provider-cloudca/issues/73))


<a name="v1.3.0"></a>
## [v1.3.0] - 2019-02-11

- Merge pull request [#69](https://github.com/cloud-ca/terraform-provider-cloudca/issues/69) from cloud-ca/dedicated_group_id
- Fixed typo
- Change go cloud-ca to 1.3.0
- Added doc
- Added dedicated_group_id field
- Update vendor dependencies ([#68](https://github.com/cloud-ca/terraform-provider-cloudca/issues/68))
- Enhance release targets and add shortcuts ([#64](https://github.com/cloud-ca/terraform-provider-cloudca/issues/64))
- Add badges to README ([#65](https://github.com/cloud-ca/terraform-provider-cloudca/issues/65))
- Add CircleCI config ([#58](https://github.com/cloud-ca/terraform-provider-cloudca/issues/58))
- Move resource document to standalone folder ([#62](https://github.com/cloud-ca/terraform-provider-cloudca/issues/62))
- Remove Glide in favor or go modules ([#57](https://github.com/cloud-ca/terraform-provider-cloudca/issues/57))
- Update year in License file ([#60](https://github.com/cloud-ca/terraform-provider-cloudca/issues/60))
- Clean up README and fix markdown lint errors ([#59](https://github.com/cloud-ca/terraform-provider-cloudca/issues/59))
- Update Offerings as of December 2018 ([#56](https://github.com/cloud-ca/terraform-provider-cloudca/issues/56))


<a name="v1.2.0"></a>
## [v1.2.0] - 2018-09-12

- Merge pull request [#51](https://github.com/cloud-ca/terraform-provider-cloudca/issues/51) from cloud-ca/development
- Updated documentation and fixed lib version for go-cloudca
- Merge pull request [#50](https://github.com/cloud-ca/terraform-provider-cloudca/issues/50) from cloud-ca/network-acl
- Updated glide lock
- Added support for network ACL name instead of ID
- Merge pull request [#49](https://github.com/cloud-ca/terraform-provider-cloudca/issues/49) from cloud-ca/ssh-keys
- Updated glide
- Updated lib version
- Added SSH key support


<a name="v1.1.0"></a>
## [v1.1.0] - 2018-04-10

- Merge pull request [#48](https://github.com/cloud-ca/terraform-provider-cloudca/issues/48) from cloud-ca/development
- Updated lib dependency
- Merge pull request [#47](https://github.com/cloud-ca/terraform-provider-cloudca/issues/47) from cloud-ca/feature/MC-4802-manual-ip-instance-creation
- updated README
- removed ForceNew flag for instance name, and return error on IP address update
- add private IP as an input
- Updated go cca lib version
- Merge pull request [#46](https://github.com/cloud-ca/terraform-provider-cloudca/issues/46) from cloud-ca/build-fix
- Merge branch 'development' of https://github.com/cloud-ca/terraform-cloudca into build-fix
- Added version to plugin
- Merge pull request [#45](https://github.com/cloud-ca/terraform-provider-cloudca/issues/45) from cloud-ca/pdion891-patch-1
- tested with terraform v0.11.1
- external providers
- the code param = root_volume_size_in_gb
- Merge pull request [#42](https://github.com/cloud-ca/terraform-provider-cloudca/issues/42) from cloud-ca/MC-5244-resize-data-volume
- Update README.md
- Modified update volume
- MC-5244: Added release-notes to make file


<a name="v1.0.2"></a>
## [v1.0.2] - 2017-05-25

- Merge branch 'development'
- Merge pull request [#41](https://github.com/cloud-ca/terraform-provider-cloudca/issues/41) from cloud-ca/MC-5244-resize-volume
- MC-5244: PR fix
- Update README.md
- MC-5244: Added release in MakeFile
- MC-5244: Updated go-cloudca deep
- Update README.md
- MC-5244: Added root volume size to instance resource
- Fixed formatting in Readme


<a name="v1.0.1"></a>
## [v1.0.1] - 2017-04-07

- Merge pull request [#39](https://github.com/cloud-ca/terraform-provider-cloudca/issues/39) from cloud-ca/upgrade-to-terraform-0.9
- [#38](https://github.com/cloud-ca/terraform-provider-cloudca/issues/38): Updated Terraform dependencies to make the provider compatible with Terraform 0.9.2


<a name="v1.0.0"></a>
## [v1.0.0] - 2017-02-22

- Replacing example of custom offerings, since none available yet
- [#36](https://github.com/cloud-ca/terraform-provider-cloudca/issues/36) Renamed tiers to networks. Updated the glide dependencies. Tested with TF version 0.8.7. Also go fit
- Merge pull request [#34](https://github.com/cloud-ca/terraform-provider-cloudca/issues/34) from cloud-ca/remove-dev-path-from-logs
- Merge pull request [#32](https://github.com/cloud-ca/terraform-provider-cloudca/issues/32) from cloud-ca/use_environment_ids_for_cca_resources
- Now trimming the $GOPATH path from source path in logs. Still some left but better than before anyways.
- Merge branch 'development' into use_environment_ids_for_cca_resources
- PR fix [#2](https://github.com/cloud-ca/terraform-provider-cloudca/issues/2)
- PR fix
- Merge pull request [#29](https://github.com/cloud-ca/terraform-provider-cloudca/issues/29) from cloud-ca/fix-packaging-checksums
- Merge pull request [#25](https://github.com/cloud-ca/terraform-provider-cloudca/issues/25) from cloud-ca/remove_purge_flag_from_instance_resource
- Merge pull request [#30](https://github.com/cloud-ca/terraform-provider-cloudca/issues/30) from cloud-ca/required_network_id_for_lbr
- Merge pull request [#31](https://github.com/cloud-ca/terraform-provider-cloudca/issues/31) from cloud-ca/handle_volume_reattach_when_instance_destroyed
- Updated doc
- Updated README
- Update README.md
- Changed the filename format of the zip file to include the name of the provider, the version and the OS/Arch. Also providing SHA256 sums of all generated files.
- Added check if current volume attached before executing detached
- Made network id required for lbw
- Handling environment not found errors
- Merge pull request [#24](https://github.com/cloud-ca/terraform-provider-cloudca/issues/24) from cloud-ca/lbr_fails_with_instance_ids
- Merge pull request [#27](https://github.com/cloud-ca/terraform-provider-cloudca/issues/27) from cloud-ca/add_private_ip_and_cidr_to_instance_and_tier
- Go fmt
- Change environment_name/service_code to environment_id
- Updated read
- Added private ip and cidr
- Removed purge flag
- Changed cast from []interface{} to *schema.Set
- Fixed network ACL rule example
- Update README.md


<a name="v0.5.0"></a>
## v0.5.0 - 2017-01-24

- Merge pull request [#15](https://github.com/cloud-ca/terraform-provider-cloudca/issues/15) from cloud-ca/development
- Updating documentation
- Changing instance ID to be required
- Update README.md
- Updated the readme to include the latest installation procedures
- Now zipping the executables. Also updated the README.md file
- Adding list + vendor to gitignore
- Changed LBR instance IDs to be a set
- Update README.md
- Update README.md
- MC-5742: Added a go dependencies system (glide) to simplify setup and reproducibility.
- Added zone to VPCs
- Merge pull request [#13](https://github.com/cloud-ca/terraform-provider-cloudca/issues/13) from cloud-ca/static-nat
- Merge pull request [#14](https://github.com/cloud-ca/terraform-provider-cloudca/issues/14) from cloud-ca/lbr
- go fmt
- go fmt
- PR feedback
- PR feedback
- Changing interface signatures to return error
- Added link in TOC
- Adding TF doc for LBR
- Add static NAT docs
- Handling stickiness policy parameters for create/update
- Move some utility functions to resource_cloudca, update static_nat resource
- Added update instances. Adding update Stickiness Method + parameters
- Adding update, still need to complete UpdateInstances + UpdateStickiness
- Fixing create and read
- Initial LBR resource need to complete CRUD
- Merge pull request [#12](https://github.com/cloud-ca/terraform-provider-cloudca/issues/12) from cloud-ca/custom_offering_in_update_instance
- Updated documentation to have custom compute offering
- Added validation to instance terraform
- Added custom compute offerings to update
- Added custom compute offering support for create instance
- Renamed compute-east to compute-qc
- Fixed an issue environment creation if the org has sub organizations and there is a username collision
- Update README.md
- Fixed doc to have change dir command
- Removed Godeps step in build process and updated doc
- Fixed typos
- Merge pull request [#11](https://github.com/cloud-ca/terraform-provider-cloudca/issues/11) from cloud-ca/remove-attach
- Implemented attach on server. No need for separate attach
- Merge pull request [#10](https://github.com/cloud-ca/terraform-provider-cloudca/issues/10) from cloud-ca/volume-custom-size-and-iops
- Update README.md
- Added support for custom size and iops in volume resource
- Renamed cloudca_publicip to cloudca_public_ip
- Update README.md
- Update README.md
- Added HCL code highlighting to doc
- Attempt to use the golang vendor feature with godep
- Updated documentation
- Merge pull request [#9](https://github.com/cloud-ca/terraform-provider-cloudca/issues/9) from cloud-ca/network-acl-rules
- MC-3708: Removed godeps folder as it was outdated. Will move to the new vendor approach.
- Update README.md
- Fixed update
- Merge pull request [#8](https://github.com/cloud-ca/terraform-provider-cloudca/issues/8) from cloud-ca/volume-resource
- Add support for zone IDs
- Fixed description
- Added some validation to optional fields
- Added missing fields to network ACL rule resource
- Update doc
- Use cloudca.Volume.GbSize instead of .Size
- Fix leftover size property references
- Renamed some go files
- Added network acl rules resource
- Change to accept integer size_in_gb, remove printf
- Merge branch 'master' into volume-resource
- Merge pull request [#7](https://github.com/cloud-ca/terraform-provider-cloudca/issues/7) from cloud-ca/Remove-name-or-id-functionality-from-resources
- Fix heading typo
- Update doc
- Implement zone ID lookup, fix attaching/detaching of volume logic
- Fixed description in tier resource
- Go fmt resource
- Merge pull request [#6](https://github.com/cloud-ca/terraform-provider-cloudca/issues/6) from cloud-ca/MC-3707-network-acl
- Removed names for ids
- Added retrieveZoneId
- formatting
- Updated to use size string
- PR fixes
- Added doc for volumes
- Update README.md
- Added Update, Delete
- Added cloudca_network_acl resource
- Added create/read in volume resource
- Added volume resource structure
- Update README.md
- Add protocol to README
- Update tier example
- Fixed VPC example
- Update README.md
- Merge pull request [#5](https://github.com/cloud-ca/terraform-provider-cloudca/issues/5) from prollynomial/master
- Add private IP id property to instance
- Remove instance_id from PFR, make end ports optional+computed
- Update resource_tier.go
- Fix inconsistent states with computed instance ids, add public+private ips as outputs
- Improve example
- Add port forwarding rule docs
- Merge remote-tracking branch 'origin/master'
- Complete create/delete/read functions
- Update README.md
- Update README.md
- Update README.md
- Added links to every resource
- Renamed Link to Links
- Added link to resources documentation
- Added service resources documentation
- Changes to public ip resource + doc
- Add port forwarding rule create
- Merge pull request [#4](https://github.com/cloud-ca/terraform-provider-cloudca/issues/4) from cloud-ca/public-ip-resource
- Moved cloudca_publicip into resources section
- Update README.md
- Update README to include public Ip
- Ran go fmt
- Additional fixes to public ip resource
- Added public IP resource
- Go fmt all files
- Merge pull request [#3](https://github.com/cloud-ca/terraform-provider-cloudca/issues/3) from cloud-ca/tier-resource
- Changed package name
- Update README.md
- Update README.md
- Update README.md
- Moving around setters
- Implemented CRUD
- Initial tier resource
- Merge branch 'vpc-resource'
- Merge pull request [#2](https://github.com/cloud-ca/terraform-provider-cloudca/issues/2) from cloud-ca/environment-resource
- Removed broken testing
- Cleaning up recurring strings
- Added environment update. Removed membership attribute
- Using sets for the users. WIP
- Missing required field
- Organization is required for env creation
- Adding value interpolation for service connection, org entry point and user to role mappings
- Adding environment create and delete. Still needs work on organization loading and role/user validation
- Added VPC resource
- Added error check to fetching of each resource in create instance
- Added insecure connection environment variable to config
- Added final step in installation instructions
- Added installation instructions + go dependencies
- Update README.md
- MC-3659: Removed optional resource data that was not symmetric with updates
- Fixed unit tests
- Fixed example
- Implemented update function and change ssh_keyname to ssh_key_name
- MC-3659: Added basic error handling
- MC-3659: Added retrieval of compute offering, template and network by Name
- MC-3659: Added retrieval of compute offering, template and network by ID
- Added description to fields of instance resource
- Fixed typo in readme
- Added How to use section in readme
- Delete terraform-provider-cloudca
- Fixed compilation issues
- Added resource_cloudca.go
- MC-3659: Move go files after review
- MC-3659: First iteration on the cloudca_instance resource
- MC-3658: Added provider skeleton and instance resource
- Added license to README
- Create LICENSE
- Added initial terraform plugin code
- Initial commit


[Unreleased]: https://github.com/cloud-ca/terraform-provider-cloudca/compare/v1.4.1...HEAD
[v1.4.1]: https://github.com/cloud-ca/terraform-provider-cloudca/compare/v1.4.0...v1.4.1
[v1.4.0]: https://github.com/cloud-ca/terraform-provider-cloudca/compare/v1.3.0...v1.4.0
[v1.3.0]: https://github.com/cloud-ca/terraform-provider-cloudca/compare/v1.2.0...v1.3.0
[v1.2.0]: https://github.com/cloud-ca/terraform-provider-cloudca/compare/v1.1.0...v1.2.0
[v1.1.0]: https://github.com/cloud-ca/terraform-provider-cloudca/compare/v1.0.2...v1.1.0
[v1.0.2]: https://github.com/cloud-ca/terraform-provider-cloudca/compare/v1.0.1...v1.0.2
[v1.0.1]: https://github.com/cloud-ca/terraform-provider-cloudca/compare/v1.0.0...v1.0.1
[v1.0.0]: https://github.com/cloud-ca/terraform-provider-cloudca/compare/v0.5.0...v1.0.0
