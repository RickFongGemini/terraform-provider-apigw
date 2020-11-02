package apigw

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type LoadBalancerCreateBody struct {
    Delay		int	`json:"delay,omitempty"`
    Desc		string	`json:"desc,omitempty"`
    ExpectedCodes	string	`json:"expected_codes,omitempty"`
    HTTPMethod		string	`json:"http_method,omitempty"`
    LBMethod		string	`json:"lb_method"`
    MaxRetries		int	`json:"max_retries,omitempty"`
    MonitorType		string	`json:"monitor_type,omitempty"`
    Name		string	`json:"name"`
    PrivateNet		string	`json:"private_net"`
    Protocol		string	`json:"protocol"`
    ProtocolPort	int	`json:"protocol_port"`
    Timeout		int	`json:"timeout,omitempty"`
    URLPath		string	`json:"url_path,omitempty"`
}

type LoadBalancerUpdateBody struct {
    LBMethod	string		`json:"lb_method,omitempty"`
    Members	*[]MemberData	`json:"members,omitempty"`
}

type MemberData struct {
    IP		string	`json:"ip,omitempty"`
    Port	int	`json:"port,omitempty"`
    Weight	int	`json:"weight,omitempty"`
}

func resourceLoadBalancer() *schema.Resource {
    return &schema.Resource{
        Create: resourceLoadBalancerCreate,
        Read:   resourceLoadBalancerRead,
        Update:	resourceLoadBalancerUpdate,
        Delete: resourceLoadBalancerDelete,

        Timeouts: &schema.ResourceTimeout{
            Create: schema.DefaultTimeout(15 * time.Minute),
            Update: schema.DefaultTimeout(15 * time.Minute),
            Delete: schema.DefaultTimeout(15 * time.Minute),
        },

        Schema: map[string]*schema.Schema{
            "active_connections": {
                Type:		schema.TypeInt,
                Computed:	true,
            },

            "create_time": {
                Type:		schema.TypeString,
                Computed:	true,
                ForceNew:	true,
            },

            "desc": {
                Type:		schema.TypeString,
                Optional:	true,
                ForceNew:	true,
            },

            "lb_method": {
                Type:		schema.TypeString,
                Required:	true,
            },

            "members": {
                Type:		schema.TypeList,
                Optional:	true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "ip": {
                            Type:	schema.TypeString,
                            Optional:	true,
                        },

                        "port": {
                            Type:	schema.TypeInt,
                            Optional:	true,
                            Default:	80,
                        },

                        "status": {
                            Type:       schema.TypeString,
                            Computed:   true,
                        },

                        "weight": {
                            Type:	schema.TypeInt,
                            Optional:	true,
                            Default:	1,
                        },
                    },
                },

                DiffSuppressFunc: lbMembersDiffFunc,
            },

            "monitor": {
                Type:		schema.TypeList,
                Optional:	true,
                Computed:	true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "delay": {
                            Type:	schema.TypeInt,
                            Optional:	true,
                            Computed:	true,
                        },

                        "expected_codes": {
                            Type:	schema.TypeString,
                            Optional:	true,
                            Computed:	true,
                        },

                        "http_method": {
                            Type:	schema.TypeString,
                            Optional:	true,
                            Computed:	true,
                        },

                        "max_retries": {
                            Type:	schema.TypeInt,
                            Optional:	true,
                            Computed:	true,
                        },

                        "monitor_type": {
                            Type:	schema.TypeString,
                            Optional:	true,
                            Computed:	true,
                        },

                        "timeout": {
                            Type:	schema.TypeInt,
                            Optional:	true,
                            Computed:	true,
                        },
                    
                        "url_path": {
                            Type:	schema.TypeString,
                            Optional:	true,
                            Computed:	true,
                        },
                    },
                },

                MaxItems:	1,
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
                Required:	true,
                ForceNew:	true,
            },

            "protocol": {
                Type:		schema.TypeString,
                Required:	true,
                ForceNew:	true,
            },

            "protocol_port": {
                Type:		schema.TypeInt,
                Required:	true,
                ForceNew:	true,
            },

            "status": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "status_reason": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "total_connections": {
                Type:		schema.TypeInt,
                Computed:	true,
            },

            "user": {       
                Type:		schema.TypeMap,
                Computed:	true,
                ForceNew:	true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },

            "vip": {
                Type:		schema.TypeString,
                Computed:	true,
            },

            "waf": {
                Type:		schema.TypeMap,
                Computed:	true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
        },
    }
}

func resourceLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    desc := d.Get("desc").(string)
    lbMethod := d.Get("lb_method").(string)
    name := d.Get("name").(string)
    platform := d.Get("platform").(string)
    privateNet := d.Get("private_net").(string)
    protocol := d.Get("protocol").(string)
    protocolPort := d.Get("protocol_port").(int)
    resourcePath := fmt.Sprintf("api/v4/%s/loadbalancers/", platform)

    body := LoadBalancerCreateBody {
        Desc:		desc,
        LBMethod:	lbMethod,
        Name:		name,
        PrivateNet:	privateNet,
        Protocol:	protocol,
        ProtocolPort:	protocolPort,
    }

    monitorArray := d.Get("monitor").([]interface{})
    if len(monitorArray) != 0 {
        info := monitorArray[0].(map[string]interface{})
        body.Delay = info["delay"].(int)
        body.ExpectedCodes = info["expected_codes"].(string)
        body.HTTPMethod = info["http_method"].(string)
        body.MaxRetries = info["max_retries"].(int)
        body.MonitorType = info["monitor_type"].(string)
        body.Timeout = info["timeout"].(int)
        body.URLPath = info["url_path"].(string)
    }

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    response, err := config.doNormalRequest(platform, resourcePath, "POST", buf)

    if err != nil {
        return fmt.Errorf("Error creating apigw_loadbalancer %s on %s: %v", name, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return err
    }

    lbID := int(data["id"].(float64))
    d.SetId(fmt.Sprintf("%d", lbID))

    newPath := fmt.Sprintf("%s/%d/", resourcePath, lbID)
    stateConf := &resource.StateChangeConf{
        Pending:    []string{"BUILD",},
        Target:     []string{"ACTIVE", "DOWN", "ERROR"},
        Refresh:    lbStateRefreshFunc(config, platform, newPath),
        Timeout:    d.Timeout(schema.TimeoutCreate),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for apigw_loadbalancer %s to become ACTIVE: %v", lbID, err)
    }

    // Update LB if user define member data
    members := d.Get("members").([]interface{})
    if len(members) > 0 {
        memberArray := make([]MemberData, len(members))
        for i, member := range members {
            detail := member.(map[string]interface{})
            memberBody := MemberData{
                IP:     detail["ip"].(string),
                Port:   detail["port"].(int),
                Weight: detail["weight"].(int),
            }

            memberArray[i] = memberBody
        }

        body := LoadBalancerUpdateBody {
            Members:	&memberArray,
        }

        buf = new(bytes.Buffer)
        json.NewEncoder(buf).Encode(body)
        _, err = config.doNormalRequest(platform, newPath, "PATCH", buf)

        if err != nil {
            return fmt.Errorf("Error updating apigw_loadbalancer %d on %s: %v", lbID, platform, err)
        }

        stateConf := &resource.StateChangeConf{
            Pending:    []string{"UPDATING"},
            Target:     []string{"ACTIVE", "ERROR"},
            Refresh:    lbStateRefreshFunc(config, platform, newPath),
            Timeout:    d.Timeout(schema.TimeoutUpdate),
            Delay:      10 * time.Second,
        }

        _, err = stateConf.WaitForState()
        if err != nil {
            return fmt.Errorf(
                "Error waiting for apigw_loadbalancer %d to become ACTIVE: %v", lbID, err)
        }
    }

    d.Set("desc", desc)
    d.Set("lb_method", lbMethod)
    d.Set("name", name)
    d.Set("platform", platform)
    d.Set("private_net", privateNet)
    d.Set("protocol", protocol)
    d.Set("protocol_port", protocolPort)
    return resourceLoadBalancerRead(d, meta)
}

func resourceLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    lbID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v4/%s/loadbalancers/%s/", platform, lbID)
    response, err := config.doNormalRequest(platform, resourcePath, "GET", nil)

    if err != nil {
        return fmt.Errorf("Unable to retrieve loadbalancer %s on %s: %v", lbID, platform, err)
    }

    var data map[string]interface{}
    err = json.Unmarshal([]byte(response), &data)

    if err != nil {
        return fmt.Errorf("Unable to retrieve loadbalancer json data: %v", err)
    }

    log.Printf("[DEBUG] Retrieved apigw_loadbalancer %s", d.Id())
    d.Set("active_connections", data["active_connections"])
    d.Set("create_time", data["create_time"])
    d.Set("lb_method", data["lb_method"])
    d.Set("members", data["members"])
    if monitor, ok := data["monitor"].(map[string]interface{}); ok {
        monitorInfo := flattenLBMonitorInfo(monitor)
        d.Set("monitor", monitorInfo)
    } else {
        d.Set("monitor", nil)
    }

    d.Set("status", data["status"])
    d.Set("status_reason", data["status_reason"])
    d.Set("total_connections", data["total_connections"])
    d.Set("user", data["user"].(map[string]interface{}))
    d.Set("vip", data["vip"])
    if waf, ok := data["waf"].(map[string]interface{}); ok {
        waf["id"] = fmt.Sprintf("%d", int(waf["id"].(float64)))
        d.Set("waf", waf)
    } else {
        d.Set("waf", data["waf"])
    }

    return nil
}

func resourceLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    var body LoadBalancerUpdateBody
    if d.HasChange("lb_method") {
        _, newLBMethod := d.GetChange("lb_method")
        body.LBMethod = newLBMethod.(string)
    }

    if d.HasChange("members") {
        _, newMembers := d.GetChange("members")
        members := newMembers.([]interface{})
        memberArray := make([]MemberData, len(members))
        for i, member := range members {
            detail := member.(map[string]interface{})
            memberBody := MemberData{
                IP:	detail["ip"].(string),
                Port:	detail["port"].(int),
                Weight:	detail["weight"].(int),
            }
             
            memberArray[i] = memberBody
        }

        body.Members = &memberArray
    }

    lbID := d.Id()
    platform := d.Get("platform").(string)
    resourcePath := fmt.Sprintf("api/v4/%s/loadbalancers/%s/", platform, lbID)

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(body)
    _, err := config.doNormalRequest(platform, resourcePath, "PATCH", buf)

    if err != nil {
        return fmt.Errorf("Error updating apigw_loadbalancer %s on %s: %v", lbID, platform, err)
    }

    stateConf := &resource.StateChangeConf{
        Pending:    []string{"UPDATING"},
        Target:     []string{"ACTIVE", "ERROR"},
        Refresh:    lbStateRefreshFunc(config, platform, resourcePath),
        Timeout:    d.Timeout(schema.TimeoutUpdate),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf(
            "Error waiting for apigw_loadbalancer %s to become ACTIVE: %v", lbID, err)
    }

    return resourceLoadBalancerRead(d, meta)
}

func resourceLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
    config := meta.(*PConfig)
    platform := d.Get("platform").(string)
    lbID := d.Id()
    resourcePath := fmt.Sprintf("api/v4/%s/loadbalancers/%s/", platform, lbID)
    _, err := config.doNormalRequest(platform, resourcePath, "DELETE", nil)

    if err != nil {
        return fmt.Errorf("Unable to delete loadbalancer %s: on %s %v", lbID, platform, err)
    }

    stateConf := &resource.StateChangeConf{
        Pending:    []string{"DELETING"},
        Target:     []string{"DELETED", "ERROR"},
        Refresh:    lbStateRefreshForDeletedFunc(config, platform, resourcePath),
        Timeout:    d.Timeout(schema.TimeoutDelete),
        Delay:      10 * time.Second,
    }

    _, err = stateConf.WaitForState()
    if err != nil {
        return fmt.Errorf( 
            "Error waiting for apigw_loadbalancer %s to become DELETED: %v", lbID, err)
    }

    d.SetId("")

    return nil
}
