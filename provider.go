// provider.go
package main

import (
        "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
  return &schema.Provider{
          DataSourcesMap: map[string]*schema.Resource{
   "gotemplate": dataSourceFile(),                  // <--- 
  },
//                ResourcesMap: map[string]*schema.Resource{
//                     "example_server": resourceServer(),
//                },
        }
}

