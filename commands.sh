# Generating swagger docs
# From root package
swag init -g cmd/devops/main.go  --parseDependency


step 1: Increment version number in code
step 2: close all milestone & related issues mentioned in it
step 3: release all binaries
# git tag -a v0.2.0 -m "m" #create dummy tag
# Execute this from the root directory where .releaser file exists
# goreleaser release --rm-dist --skip-validate (Run this for core binary as well as for k8s plugin from that directory)
step 4: Create a new release on devops-cli repository, by mentioning the milestone
step 5: Delete issues from github project
