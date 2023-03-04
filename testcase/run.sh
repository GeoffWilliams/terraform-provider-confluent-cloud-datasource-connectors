pushd ..
make build
popd
rm -f .terraform -rf
rm -f .terraform.lock.hcl
terraform init
TF_LOG=DEBUG terraform apply -auto-approve
terraform output