step 1: Increment version number in code
step 2: close all milestone & related issues mentioned in it
step 3: release all binaries
# goreleaser release --rm-dist --skip-validate
step 4: Create a new release on devops-cli repository, by mentioning the milestone
step 5: Delete issues from github project
