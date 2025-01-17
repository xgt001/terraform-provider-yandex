package yandex

import (
	"bytes"
	"fmt"
	"net"
	"strings"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/loadbalancer/v1"
)

func resourceLBTargetGroupTargetHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["subnet_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["address"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func resourceLBNetowrkLoadBalancerListenerHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func resourceLBNetowrkLoadBalancerAttachedTargetGroupHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["target_group_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	switch hcs := m["healthcheck"].(type) {
	case []interface{}:
		for _, hcc := range hcs {
			hc := hcc.(map[string]interface{})
			if v, ok := hc["name"]; ok {
				buf.WriteString(fmt.Sprintf("%s-", v.(string)))
			}
		}
	case []map[string]interface{}:
		for _, hc := range hcs {
			if v, ok := hc["name"]; ok {
				buf.WriteString(fmt.Sprintf("%s-", v.(string)))
			}
		}
	default:
	}

	return hashcode.String(buf.String())
}

func expandLBListenerSpecs(d *schema.ResourceData) ([]*loadbalancer.ListenerSpec, error) {
	var result []*loadbalancer.ListenerSpec
	listenersSet := d.Get("listener").(*schema.Set)

	for _, v := range listenersSet.List() {
		config := v.(map[string]interface{})

		ls, err := expandLBListenerSpec(config)
		if err != nil {
			return nil, err
		}

		result = append(result, ls)
	}

	return result, nil
}

func expandLBListenerSpec(config map[string]interface{}) (*loadbalancer.ListenerSpec, error) {
	ls := &loadbalancer.ListenerSpec{}

	if v, ok := config["name"]; ok {
		ls.Name = v.(string)
	}

	if v, ok := config["port"]; ok {
		ls.Port = int64(v.(int))
	}

	if v, ok := config["target_port"]; ok {
		ls.TargetPort = int64(v.(int))
	}

	if v, ok := config["protocol"]; ok {
		p, err := parseListenerProtocol(v.(string))
		if err != nil {
			return nil, err
		}
		ls.Protocol = p
	}

	if v, ok := getFirstElement(config, "external_address_spec"); ok {
		eas, err := expandLBExternalAddressSpec(v)
		if err != nil {
			return nil, err
		}
		ls.Address = eas
	}

	return ls, nil
}

func expandLBExternalAddressSpec(config map[string]interface{}) (*loadbalancer.ListenerSpec_ExternalAddressSpec, error) {
	as := &loadbalancer.ListenerSpec_ExternalAddressSpec{
		ExternalAddressSpec: &loadbalancer.ExternalAddressSpec{},
	}

	if v, ok := config["address"]; ok {
		as.ExternalAddressSpec.Address = v.(string)
	}

	if v, ok := config["ip_version"]; ok {
		v, err := parseExternalIPVersion(v.(string))
		if err != nil {
			return nil, err
		}
		as.ExternalAddressSpec.IpVersion = v
	}

	return as, nil
}

func expandLBAttachedTargetGroups(d *schema.ResourceData) ([]*loadbalancer.AttachedTargetGroup, error) {
	var result []*loadbalancer.AttachedTargetGroup
	atgsSet := d.Get("attached_target_group").(*schema.Set)

	for _, v := range atgsSet.List() {
		config := v.(map[string]interface{})

		atg, err := expandLBAttachedTargetGroup(config)
		if err != nil {
			return nil, err
		}

		result = append(result, atg)
	}

	return result, nil
}

func expandLBAttachedTargetGroup(config map[string]interface{}) (*loadbalancer.AttachedTargetGroup, error) {
	atg := &loadbalancer.AttachedTargetGroup{}

	if v, ok := config["target_group_id"]; ok {
		atg.TargetGroupId = v.(string)
	}

	if v, ok := config["healthcheck"]; ok {
		hcConfigs := v.([]interface{})
		atg.HealthChecks = make([]*loadbalancer.HealthCheck, len(hcConfigs))

		for i := 0; i < len(hcConfigs); i++ {
			hcConfig := hcConfigs[i]
			hc, err := expandLBHealthcheck(hcConfig.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			atg.HealthChecks[i] = hc
		}
	}

	return atg, nil
}

func expandLBHealthcheck(config map[string]interface{}) (*loadbalancer.HealthCheck, error) {
	hc := &loadbalancer.HealthCheck{}

	if v, ok := config["name"]; ok {
		hc.Name = v.(string)
	}

	if v, ok := config["interval"]; ok {
		hc.Interval = &duration.Duration{Seconds: int64(v.(int))}
	}

	if v, ok := config["timeout"]; ok {
		hc.Timeout = &duration.Duration{Seconds: int64(v.(int))}
	}

	if v, ok := config["unhealthy_threshold"]; ok {
		hc.UnhealthyThreshold = int64(v.(int))
	}

	if v, ok := config["healthy_threshold"]; ok {
		hc.HealthyThreshold = int64(v.(int))
	}

	httpOptions, httpOptionsOk := getFirstElement(config, "http_options")
	tcpOptions, tcpOptionsOk := getFirstElement(config, "tcp_options")

	if httpOptionsOk && tcpOptionsOk {
		return nil, fmt.Errorf("Use one of 'http_options' or 'tcp_options', not both")
	}

	if httpOptionsOk {
		options, err := expandLBHealthcheckHTTPOptions(httpOptions)
		if err != nil {
			return nil, err
		}
		hc.Options = &loadbalancer.HealthCheck_HttpOptions_{
			HttpOptions: options,
		}
	}

	if tcpOptionsOk {
		options, err := expandLBHealthcheckTCPOptions(tcpOptions)
		if err != nil {
			return nil, err
		}
		hc.Options = &loadbalancer.HealthCheck_TcpOptions_{
			TcpOptions: options,
		}
	}

	return hc, nil
}

func expandLBHealthcheckHTTPOptions(config map[string]interface{}) (*loadbalancer.HealthCheck_HttpOptions, error) {
	options := &loadbalancer.HealthCheck_HttpOptions{}

	if v, ok := config["port"]; ok {
		options.Port = int64(v.(int))
	}

	if v, ok := config["path"]; ok {
		options.Path = v.(string)
	}

	return options, nil
}

func expandLBHealthcheckTCPOptions(config map[string]interface{}) (*loadbalancer.HealthCheck_TcpOptions, error) {
	options := &loadbalancer.HealthCheck_TcpOptions{}

	if v, ok := config["port"]; ok {
		options.Port = int64(v.(int))
	}

	return options, nil
}

func expandLBTargets(d *schema.ResourceData) ([]*loadbalancer.Target, error) {
	var targets []*loadbalancer.Target
	targetsSet := d.Get("target").(*schema.Set)

	for _, t := range targetsSet.List() {
		targetConfig := t.(map[string]interface{})

		target, err := expandLBTarget(targetConfig)
		if err != nil {
			return nil, err
		}

		targets = append(targets, target)
	}

	return targets, nil
}

func expandLBTarget(config map[string]interface{}) (*loadbalancer.Target, error) {
	target := &loadbalancer.Target{}

	if v, ok := config["subnet_id"]; ok {
		target.SubnetId = v.(string)
	}

	if v, ok := config["address"]; ok {
		target.Address = v.(string)
	}

	return target, nil
}

func flattenLBTargets(tg *loadbalancer.TargetGroup) (*schema.Set, error) {
	result := &schema.Set{F: resourceLBTargetGroupTargetHash}

	for _, t := range tg.Targets {
		flTarget := map[string]interface{}{
			"subnet_id": t.SubnetId,
			"address":   t.Address,
		}
		result.Add(flTarget)
	}

	return result, nil
}

func flattenLBListenerSpecs(nlb *loadbalancer.NetworkLoadBalancer) (*schema.Set, error) {
	result := &schema.Set{F: resourceLBNetowrkLoadBalancerListenerHash}

	for _, ls := range nlb.Listeners {
		eas, err := flattenLBExternalAddressSpec(ls)
		if err != nil {
			return nil, err
		}
		flListener := map[string]interface{}{
			"name":                  ls.Name,
			"port":                  ls.Port,
			"target_port":           ls.TargetPort,
			"protocol":              strings.ToLower(ls.Protocol.String()),
			"external_address_spec": []map[string]interface{}{eas},
		}
		result.Add(flListener)
	}

	return result, nil
}

func flattenLBExternalAddressSpec(ls *loadbalancer.Listener) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"address": ls.Address,
	}

	addr := net.ParseIP(ls.Address)
	isV4 := addr.To4() != nil
	if isV4 {
		result["ip_version"] = "ipv4"
	} else {
		result["ip_version"] = "ipv6"
	}

	return result, nil
}

func flattenLBAttachedTargetGroups(nlb *loadbalancer.NetworkLoadBalancer) (*schema.Set, error) {
	result := &schema.Set{F: resourceLBNetowrkLoadBalancerAttachedTargetGroupHash}

	for _, atg := range nlb.AttachedTargetGroups {
		hcs, err := flattenLBHealthchecks(atg)
		if err != nil {
			return nil, err
		}

		flATG := map[string]interface{}{
			"target_group_id": atg.TargetGroupId,
			"healthcheck":     hcs,
		}
		result.Add(flATG)
	}

	return result, nil
}

func flattenLBHealthchecks(atg *loadbalancer.AttachedTargetGroup) ([]map[string]interface{}, error) {
	result := []map[string]interface{}{}

	for _, hc := range atg.HealthChecks {
		flHC := map[string]interface{}{
			"name":                hc.Name,
			"interval":            hc.Interval.Seconds,
			"timeout":             hc.Timeout.Seconds,
			"unhealthy_threshold": hc.UnhealthyThreshold,
			"healthy_threshold":   hc.HealthyThreshold,
		}
		switch hc.Options.(type) {
		case *loadbalancer.HealthCheck_HttpOptions_:
			flHC["http_options"] = []map[string]interface{}{
				{
					"port": hc.GetHttpOptions().Port,
					"path": hc.GetHttpOptions().Path,
				},
			}
		case *loadbalancer.HealthCheck_TcpOptions_:
			flHC["tcp_options"] = []map[string]interface{}{
				{
					"port": hc.GetTcpOptions().Port,
				},
			}
		default:
			return nil, fmt.Errorf("Unknown healthcheck options type: %T", hc.Options)
		}
		result = append(result, flHC)
	}

	return result, nil
}

func parseListenerProtocol(s string) (loadbalancer.Listener_Protocol, error) {
	switch strings.ToLower(s) {
	case "tcp":
		return loadbalancer.Listener_TCP, nil
	case "":
		return loadbalancer.Listener_TCP, nil
	default:
		return loadbalancer.Listener_PROTOCOL_UNSPECIFIED,
			fmt.Errorf("value for 'protocol' must be 'tcp', not '%s'", s)
	}
}

func parseNetworkLoadBalancerType(s string) (loadbalancer.NetworkLoadBalancer_Type, error) {
	switch s {
	case "external":
		return loadbalancer.NetworkLoadBalancer_EXTERNAL, nil
	case "":
		return loadbalancer.NetworkLoadBalancer_EXTERNAL, nil
	default:
		return loadbalancer.NetworkLoadBalancer_TYPE_UNSPECIFIED,
			fmt.Errorf("value for 'type' must be 'external', not '%s'", s)
	}
}

func parseExternalIPVersion(s string) (loadbalancer.IpVersion, error) {
	switch strings.ToLower(s) {
	case "ipv4":
		return loadbalancer.IpVersion_IPV4, nil
	case "ipv6":
		return loadbalancer.IpVersion_IPV6, nil
	case "":
		return loadbalancer.IpVersion_IPV4, nil
	default:
		return loadbalancer.IpVersion_IP_VERSION_UNSPECIFIED,
			fmt.Errorf("value for 'external ip version' must be 'ipv4' or 'ipv6', not '%s'", s)
	}
}

func getFirstElement(config map[string]interface{}, name string) (map[string]interface{}, bool) {
	if v, ok := config[name]; ok {
		cfgList := v.([]interface{})
		if len(cfgList) > 0 {
			return cfgList[0].(map[string]interface{}), true
		}
	}
	return nil, false
}
