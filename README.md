# terraform-provider-confluent-cloud-datasource-connectors

This is an experiment to read data from the Confluent Cloud API to determine whether a named connector exists or not.

Useful output is an argument `.status` containing either `DEFINED` or `NOT_DEFINED`

Project is heavily derived from [confluent-terraform-provider](https://docs.confluent.io/cloud/current/get-started/terraform-provider.html#using-the-confluent-terraform-provider)

## Example Usage

see [testcase](testcase/test.tf)


# deploy plugin locally for dev/test

https://stackoverflow.com/questions/68182628/terraform-use-local-provider-plugin

## Terraform rc
```
cat <<EOF > ~/.terraformrc
provider_installation {
  filesystem_mirror {
    path    = "${HOME}/.terraform.d/plugins"
  }
  direct {
    exclude = ["terraform.local/*/*"]
  }
}
EOF
```

## Directory for the plugin
```
mkdir -p ~/.terraform.d/plugins/terraform.local/local/confluent-cloud-datasource-connectors/0.0.1/linux_amd64/
```

## Symlink into place
```
ln -fs $(pwd)/terraform-provider-confluent-cloud-datasource-connectors_v0.0.1 ~/.terraform.d/plugins/terraform.local/local/confluent-cloud-datasource-connectors/0.0.1/linux_amd64/terraform-provider-confluent-cloud-datasource-connectors_v0.0.1
```

