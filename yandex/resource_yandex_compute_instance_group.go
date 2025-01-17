package yandex

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1/instancegroup"
)

const (
	yandexComputeInstanceGroupDefaultTimeout = 30 * time.Minute
)

func resourceYandexComputeInstanceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceYandexComputeInstanceGroupCreate,
		Read:   resourceYandexComputeInstanceGroupRead,
		Update: resourceYandexComputeInstanceGroupUpdate,
		Delete: resourceYandexComputeInstanceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(yandexComputeInstanceGroupDefaultTimeout),
			Update: schema.DefaultTimeout(yandexComputeInstanceGroupDefaultTimeout),
			Delete: schema.DefaultTimeout(yandexComputeInstanceGroupDefaultTimeout),
		},

		SchemaVersion: 0,

		Schema: map[string]*schema.Schema{
			"service_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_template": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resources": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"memory": {
										Type:         schema.TypeFloat,
										Required:     true,
										ValidateFunc: FloatAtLeast(0.0),
									},

									"cores": {
										Type:     schema.TypeInt,
										Required: true,
									},

									"gpus": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},

									"core_fraction": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  100,
									},
								},
							},
						},

						"boot_disk": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"initialize_params": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"description": {
													Type:     schema.TypeString,
													Optional: true,
												},

												"size": {
													Type:         schema.TypeInt,
													Optional:     true,
													Computed:     true,
													ValidateFunc: validation.IntAtLeast(1),
												},

												"type": {
													Type:         schema.TypeString,
													Optional:     true,
													Default:      "network-hdd",
													ValidateFunc: validation.StringInSlice([]string{"network-hdd", "network-nvme"}, false),
												},

												"image_id": {
													Type:          schema.TypeString,
													Optional:      true,
													Computed:      true,
													ConflictsWith: []string{"instance_template.0.boot_disk.initialize_params.snapshot_id"},
												},

												"snapshot_id": {
													Type:          schema.TypeString,
													Optional:      true,
													Computed:      true,
													ConflictsWith: []string{"instance_template.0.boot_disk.initialize_params.image_id"},
												},
											},
										},
									},

									"mode": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "READ_WRITE",
										ValidateFunc: validation.StringInSlice([]string{"READ_WRITE"}, false),
									},

									"device_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},

						"network_interface": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"network_id": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"subnet_ids": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"nat": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},

									"ipv6": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},

						"platform_id": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "standard-v1",
						},

						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"metadata": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"labels": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"secondary_disk": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"initialize_params": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"description": {
													Type:     schema.TypeString,
													Optional: true,
												},

												"size": {
													Type:         schema.TypeInt,
													Optional:     true,
													ValidateFunc: validation.IntAtLeast(1),
													Default:      8,
												},

												"type": {
													Type:         schema.TypeString,
													Optional:     true,
													Default:      "network-hdd",
													ValidateFunc: validation.StringInSlice([]string{"network-hdd", "network-nvme"}, false),
												},

												"image_id": {
													Type:          schema.TypeString,
													Optional:      true,
													ConflictsWith: []string{"instance_template.secondary_disk.initialize_params.snapshot_id"},
												},

												"snapshot_id": {
													Type:          schema.TypeString,
													Optional:      true,
													ConflictsWith: []string{"instance_template.secondary_disk.initialize_params.image_id"},
												},
											},
										},
									},

									"mode": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "READ_WRITE",
										ValidateFunc: validation.StringInSlice([]string{"READ_ONLY", "READ_WRITE"}, false),
									},

									"device_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"scheduling_policy": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"preemptible": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},

						"service_account_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"scale_policy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fixed_scale": {
							Type:          schema.TypeList,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: []string{"scale_policy.0.auto_scale"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
						"auto_scale": {
							Type:          schema.TypeList,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: []string{"scale_policy.0.fixed_scale"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"initial_size": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"measurement_duration": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"cpu_utilization_target": {
										Type:     schema.TypeFloat,
										Required: true,
									},
									"min_zone_size": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"max_size": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"warmup_duration": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"stabilization_duration": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"deploy_policy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_unavailable": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"max_expansion": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"max_deleting": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"max_creating": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"startup_duration": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},

			"allocation_policy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zones": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"folder_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"health_check": {
				Type:     schema.TypeList,
				MinItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"healthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  2,
						},

						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  2,
						},

						"tcp_options": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"health_check.http_options"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},

						"http_options": {
							Type:          schema.TypeList,
							Optional:      true,
							ConflictsWith: []string{"health_check.tcp_options"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port": {
										Type:     schema.TypeInt,
										Required: true,
									},

									"path": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"load_balancer": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_group_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"target_group_description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"target_group_labels": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"target_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceYandexComputeInstanceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	req, err := prepareCreateInstanceGroupRequest(d, config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(config.Context(), d.Timeout(schema.TimeoutCreate))
	defer cancel()

	op, err := config.sdk.WrapOperation(config.sdk.InstanceGroup().InstanceGroup().Create(ctx, req))
	if err != nil {
		return fmt.Errorf("Error while requesting API to create instance group: %s", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Error while waiting operation to create instance group: %s", err)
	}

	resp, err := op.Response()
	if err != nil {
		return fmt.Errorf("Instance group creation failed: %s", err)
	}

	instanceGroup, ok := resp.(*instancegroup.InstanceGroup)
	if !ok {
		return fmt.Errorf("Create response doesn't contain Instance group")
	}

	d.SetId(instanceGroup.Id)

	return resourceYandexComputeInstanceGroupRead(d, meta)
}

func resourceYandexComputeInstanceGroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ctx, cancel := context.WithTimeout(config.Context(), d.Timeout(schema.TimeoutRead))
	defer cancel()

	instanceGroup, err := config.sdk.InstanceGroup().InstanceGroup().Get(ctx, &instancegroup.GetInstanceGroupRequest{
		InstanceGroupId: d.Id(),
		View:            instancegroup.InstanceGroupView_FULL,
	})

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Instance group %q", d.Id()))
	}

	return flattenInstanceGroup(d, instanceGroup)
}

func flattenInstanceGroup(d *schema.ResourceData, instanceGroup *instancegroup.InstanceGroup) error {
	createdAt, err := getTimestamp(instanceGroup.CreatedAt)
	if err != nil {
		return err
	}

	d.Set("created_at", createdAt)
	d.Set("folder_id", instanceGroup.GetFolderId())
	d.Set("name", instanceGroup.GetName())
	d.Set("description", instanceGroup.GetDescription())
	d.Set("service_account_id", instanceGroup.GetServiceAccountId())

	if err := d.Set("labels", instanceGroup.GetLabels()); err != nil {
		return err
	}

	template, err := flattenInstanceGroupInstanceTemplate(instanceGroup.GetInstanceTemplate())
	if err != nil {
		return err
	}
	if err := d.Set("instance_template", template); err != nil {
		return err
	}

	scalePolicy, err := flattenInstanceGroupScalePolicy(instanceGroup)
	if err != nil {
		return err
	}
	if err := d.Set("scale_policy", scalePolicy); err != nil {
		return err
	}

	deployPolicy, err := flattenInstanceGroupDeployPolicy(instanceGroup)
	if err != nil {
		return err
	}
	if err := d.Set("deploy_policy", deployPolicy); err != nil {
		return err
	}

	allocationPolicy, err := flattenInstanceGroupAllocationPolicy(instanceGroup)
	if err != nil {
		return err
	}

	if err := d.Set("allocation_policy", allocationPolicy); err != nil {
		return err
	}

	loadBalancerSpec, err := flattenInstanceGroupLoadBalancerSpec(instanceGroup)
	if err != nil {
		return err
	}

	if err := d.Set("load_balancer", loadBalancerSpec); err != nil {
		return err
	}

	healthChecks, err := flattenInstanceGroupHealthChecks(instanceGroup)
	if err != nil {
		return err
	}

	return d.Set("health_check", healthChecks)
}

func resourceYandexComputeInstanceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	req, err := prepareUpdateInstanceGroupRequest(d, config)
	if err != nil {
		return err
	}

	err = makeInstanceGroupUpdateRequest(req, d, meta)
	if err != nil {
		return err
	}

	return resourceYandexComputeInstanceGroupRead(d, meta)
}

func resourceYandexComputeInstanceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	log.Printf("[DEBUG] Deleting Instance group %q", d.Id())

	req := &instancegroup.DeleteInstanceGroupRequest{
		InstanceGroupId: d.Id(),
	}

	ctx, cancel := context.WithTimeout(config.Context(), d.Timeout(schema.TimeoutDelete))
	defer cancel()

	op, err := config.sdk.WrapOperation(config.sdk.InstanceGroup().InstanceGroup().Delete(ctx, req))
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Instance group %q", d.Get("name").(string)))
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	_, err = op.Response()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Instance group %q", d.Id())
	return nil
}

func prepareCreateInstanceGroupRequest(d *schema.ResourceData, meta *Config) (*instancegroup.CreateInstanceGroupRequest, error) {
	folderID, err := getFolderID(d, meta)
	if err != nil {
		return nil, fmt.Errorf("Error getting folder ID while creating instance group: %s", err)
	}

	labels, err := expandLabels(d.Get("labels"))
	if err != nil {
		return nil, fmt.Errorf("Error expanding labels while creating instance group: %s", err)
	}

	instanceTemplate, err := expandInstanceGroupInstanceTemplate(d, "instance_template.0", meta)
	if err != nil {
		return nil, fmt.Errorf("Error creating 'instance_template' object of api request: %s", err)
	}

	scalePolicy, err := expandInstanceGroupScalePolicy(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating 'scale_policy' object of api request: %s", err)
	}

	deployPolicy, err := expandInstanceGroupDeployPolicy(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating 'deploy_policy' object of api request: %s", err)
	}

	allocationPolicy, err := expandInstanceGroupAllocationPolicy(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating 'allocation_policy' object of api request: %s", err)
	}

	healthChecksSpec, err := expandInstanceGroupHealthCheckSpec(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating 'health_checks_spec' object of api request: %s", err)
	}

	loadBalancerSpec, err := expandInstanceGroupLoadBalancerSpec(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating 'load_balancer_spec' object of api request: %s", err)
	}

	req := &instancegroup.CreateInstanceGroupRequest{
		FolderId:         folderID,
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Labels:           labels,
		InstanceTemplate: instanceTemplate,
		ScalePolicy:      scalePolicy,
		DeployPolicy:     deployPolicy,
		AllocationPolicy: allocationPolicy,
		LoadBalancerSpec: loadBalancerSpec,
		HealthChecksSpec: healthChecksSpec,
		ServiceAccountId: d.Get("service_account_id").(string),
	}

	return req, nil
}

func prepareUpdateInstanceGroupRequest(d *schema.ResourceData, meta *Config) (*instancegroup.UpdateInstanceGroupRequest, error) {
	labels, err := expandLabels(d.Get("labels"))
	if err != nil {
		return nil, fmt.Errorf("Error expanding labels while creating instance: %s", err)
	}

	instanceTemplate, err := expandInstanceGroupInstanceTemplate(d, "instance_template.0", meta)
	if err != nil {
		return nil, fmt.Errorf("Error creating 'instance_template' object of api request: %s", err)
	}

	scalePolicy, err := expandInstanceGroupScalePolicy(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating 'scale_policy' object of api request: %s", err)
	}

	deployPolicy, err := expandInstanceGroupDeployPolicy(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating 'deploy_policy' object of api request: %s", err)
	}

	allocationPolicy, err := expandInstanceGroupAllocationPolicy(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating 'allocation_policy' object of api request: %s", err)
	}

	healthChecksSpec, err := expandInstanceGroupHealthCheckSpec(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating 'health_checks_spec' object of api request: %s", err)
	}

	loadBalancerSpec, err := expandInstanceGroupLoadBalancerSpec(d)
	if err != nil {
		return nil, fmt.Errorf("Error creating 'load_balancer_spec' object of api request: %s", err)
	}

	req := &instancegroup.UpdateInstanceGroupRequest{
		InstanceGroupId:  d.Id(),
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Labels:           labels,
		InstanceTemplate: instanceTemplate,
		ScalePolicy:      scalePolicy,
		DeployPolicy:     deployPolicy,
		AllocationPolicy: allocationPolicy,
		LoadBalancerSpec: loadBalancerSpec,
		HealthChecksSpec: healthChecksSpec,
		ServiceAccountId: d.Get("service_account_id").(string),
	}

	return req, nil
}

func makeInstanceGroupUpdateRequest(req *instancegroup.UpdateInstanceGroupRequest, d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	ctx, cancel := context.WithTimeout(config.Context(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()

	op, err := config.sdk.WrapOperation(config.sdk.InstanceGroup().InstanceGroup().Update(ctx, req))
	if err != nil {
		return fmt.Errorf("Error while requesting API to update Instance group %q: %s", d.Id(), err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Error updating Instance group %q: %s", d.Id(), err)
	}

	return nil
}
