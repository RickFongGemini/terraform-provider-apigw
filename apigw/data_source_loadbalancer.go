package apigw

import (
    "fmt"
    "log"
    "encoding/json"
    "strings"

    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceLoadBalancer() * schema.Resource {
    return &schema.Resource{
        Read: dataSourceLoadBalancerRead,

        Schema: map[string]*schema.Schema{
            "create_time": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "desc": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "lb_method": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "name": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "platform": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "private_net": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "project": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "protocol": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "protocol_port": {
                Type:		schema.TypeInt,
                Computed:	true,
            },

            "status": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "user": {
                Type:		schema.TypeMap,
                Computed:	true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
        },
    }
}

// dataSourceLoadBalancerRead performs the loadbalancer lookup.
func dataSourceLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)

    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    params := []string{fmt.Sprintf("name=%s", name)}
    if project := d.Get("project"); project != "" {
        params = append(params, fmt.Sprintf("project=%s", project))
    }
    if private_net := d.Get("private_net"); private_net != "" {
        params = append(params, fmt.Sprintf("private_net=%s", private_net))
    }
    if len(params) < 2 {
        return fmt.Errorf("Either project or private_net should be defined")
    }

    resourcePath := fmt.Sprintf("api/v4/%s/loadbalancers/?%s", platform,
                                strings.Join(params, "&"))
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to list loadbalancers: %v", err)
    }

    var data []map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    for _, loadbalancer := range data {
        if loadbalancer["name"] == name {
            return dataSourceLoadBalancerAttributes(d, loadbalancer)
        }
    }

    return fmt.Errorf("Unable to retrieve loadbalancer %s: %v", name, err)
}

// dataSourceLoadBalancerAttributes populates the fields of a loadbalancer data source.
func dataSourceLoadBalancerAttributes(d *schema.ResourceData, data map[string]interface{}) error {
    loadbalancer_id := int(data["id"].(float64))
    log.Printf("[DEBUG] Retrieved apigw_loadbalancer: %d", loadbalancer_id)

    d.SetId(fmt.Sprintf("%d", loadbalancer_id))
    d.Set("status", data["status"])
    d.Set("user", data["user"])
    d.Set("name", data["name"])
    d.Set("desc", data["desc"])
    d.Set("protocol", data["protocol"])
    d.Set("protocol_port", data["protocol_port"])
    d.Set("lb_method", data["lb_method"])
    private_net_info := data["private_net"].(map[string]interface{})
    private_net := fmt.Sprintf("%d", int(private_net_info["id"].(float64)))
    d.Set("private_net", private_net)

    return nil
}
