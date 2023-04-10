# Command to generate initial code from swagger.yaml
openapi-generator-cli generate  -i ./swagger.yaml -o src/generated-sources/openapi -g typescript-fetch --additional-properties=supportsES6=true,npmVersion=6.9.0,typescriptThreePlus=true

# For building the project
# Uncomment the production URL
vite build
tar -czvf ui.tar.gz dist/
cp ui.tar.gz ~/.devops/
tar -xvzf ui.tar.gz