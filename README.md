## Example Usage

```hcl
terraform {
  required_providers {
    env = {
      source  = "clockworksoul/env"
      version = "0.0.2"
    }
  }
}
```

## Using the Value Data Source

```hcl
data "env_value" "environment" {
  key = "ENV"
}

resource "aws_instance" "web" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t3.micro"

  tags = {
    Name = "HelloWorld"
    env  = data.env_value.environment.value
  }
}
```

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

