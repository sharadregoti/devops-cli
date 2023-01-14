step 1: Increment version number in code
step 2: close all milestone & related issues mentioned in it
step 3: release all binaries
# Change version in install2.sh & code 
# git tag -a v0.2.0 -m "m" #create dummy tag
# Execute this from the root directory where .releaser file exists
# goreleaser release --rm-dist --skip-validate
step 4: Create a new release on devops-cli repository, by mentioning the milestone
step 5: Delete issues from github project
