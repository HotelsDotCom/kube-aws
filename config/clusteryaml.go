package config

import (
	"github.com/kubernetes-incubator/kube-aws-ng/model"
	"github.com/kubernetes-incubator/kube-aws-ng/types"
	"github.com/kubernetes-incubator/kube-aws-ng/types/coreos"
	"github.com/kubernetes-incubator/kube-aws-ng/types/dex"
	"github.com/kubernetes-incubator/kube-aws-ng/types/ec2"
	"github.com/kubernetes-incubator/kube-aws-ng/types/kubernetes"
)

type SubnetConfRef struct {
	Name ec2.SubnetName
}

type SubnetIdConf struct {
	Id                ec2.SubnetId  //conflicts: IdFromStackOutput
	IdFromStackOutput ec2.StackName //conflicts: Id;
}

type NGWIdConf struct {
	Id                ec2.NGWId      //conflicts: IdFromStackOutput
	IdFromStackOutput ec2.StackName  //conflicts: Id;
	eipAllocationId   ec2.EIPAllocId //conflicts: Id, IdFromStackOutput
}

type RouteTableIdConf struct {
	Id                ec2.RouteTableId //conflicts: IdFromStackOutput
	IdFromStackOutput ec2.StackName    //conflicts: Id;
}

type SubnetConf struct {
	Name             ec2.SubnetName
	AvailabilityZone ec2.AvailabilityZone `yaml:"availabilityZone"`
	InstanceCIDR     types.IPNet `yaml:"instanceCIDR"`
	Private          bool
	SubnetIdConf `yaml:",inline"`
	NatGateway NGWIdConf
	RouteTable RouteTableIdConf
}

type APIEndpointLoadBalancer struct { // need invariants generator
	Id              ec2.ELBName
	CreateRecordSet bool  `yaml:"createRecordSet"`
	RecordSetTTL    uint `yaml:"recordSetTTL"`

	Subnets                     []SubnetConfRef
	Private                     bool
	HostedZone                  ec2.HostedZoneId `yaml:"hostedZone"`
	ApiAccessAllowedSourceCIDRs []types.IPNet `yaml:"ApiAccessAllowedSourceCIDRs"`
	SecurityGroupIds            []ec2.SecurityGroupId `yaml:"securityGroupIds"`
}

type APIEndpointName string

type APIEndpointConf struct {
	Name    APIEndpointName
	DnsName types.DNSName `yaml:"dnsName"`

	LoadBalancer APIEndpointLoadBalancer `yaml:"loadBalancer"`
}

type ASGConf struct {
	MinSize                            uint `yaml:"minSize"`
	MaxSize                            uint `yaml:"maxSize"`
	RollingUpdateMinInstancesInService uint `yaml:"rollingUpdateMinInstancesInService"`
}

type IAMConf struct {
	Role struct {
		ManagedPolicies []ec2.IAMPolicyARN  `yaml:"managedPolicies"`
		InstanceProfile ec2.InstanceProfileARN `yaml:"instanceProfile"`
	}
}

type VolumeConf struct {
	Size uint
	Type ec2.VolumeType
	Iops uint
}

type VolumeMountConf struct {
	VolumeConf
	Device types.BlockDeviceName
	Path   types.FilesystemPath
}

type MaybeEncryptedOrEphemeralVolume struct {
	VolumeConf
	Encrypted bool `yaml:encrypted` // conflicts: Ephemeral
	Ephemeral bool `yaml:ephemeral` // conflicts: Encrypted, broken
}

type InstanceCommonDescrEmbed struct {
	Count              uint 
	CreateTimeout      ec2.Timeout `yaml:"createTimeout"`
	InstanceType       ec2.InstanceType `yaml:"instanceType"`
	Tenancy                   ec2.InstanceTenancy //new for Controller
	RootVolume         VolumeConf `yaml:rootVolume` //validate: must be empty/default device and path
	SecurityGroupIds   []ec2.SecurityGroupId `yaml:SecurityGroupIds`
	IAM                IAMConf
	Subnets            []SubnetConfRef
	KeyName        ec2.SSHKeyPairName  // new for Etcd,Controller
	ReleaseChannel coreos.ReleaseChannel //  new for Etcd,Controller
	AmiId          ec2.AmiId   // new for Etcd, Controller
	ManagedIamRoleSuffix      ec2.IAMRoleName `yaml:mangedIamRoleName`  //new for Etcd,Ctrl
	CustomFiles        []map[string]interface{}
	CustomSystemdUnits []map[string]interface{}
}

type ControllerConf struct {
	InstanceCommonDescrEmbed `yaml:",inline"`
	NodeLabels         map[kubernetes.LabelName]kubernetes.LabelValue `yaml:"nodeLabels"`
	AutoScalingGroup   ASGConf `yaml:"autoScalingGroup"`
	LoadBalancer       struct {
		Private bool
		Subnets []SubnetConfRef
	}       `yaml:"loadBalancer,omitempty"` // how is it connected to apiEndpoints ELBs?
}

// validate: g2 or p2 instance, docker runtime
type GpuConf struct {
	Nvidia struct {
		Enabled bool
		Version types.NvidiaDriverVersion
	}
}

type NodePoolName string

type SpotFleetConf struct{} //TODO

type ContainerImages struct {
        HyperkubeImage                     model.Image `yaml:"hyperkubeImage,omitempty",default:"{'quay.io/coreos/hyperkube', 'v1.7.3_coreos.0', false}"`
        AWSCliImage                        model.Image `yaml:"awsCliImage,omitempty"`
        CalicoNodeImage                    model.Image `yaml:"calicoNodeImage,omitempty"`
        CalicoCniImage                     model.Image `yaml:"calicoCniImage,omitempty"`
        CalicoCtlImage                     model.Image `yaml:"calicoCtlImage,omitempty"`
        CalicoPolicyControllerImage        model.Image `yaml:"calicoPolicyControllerImage,omitempty"`
        ClusterAutoscalerImage             model.Image `yaml:"clusterAutoscalerImage,omitempty"`
        ClusterProportionalAutoscalerImage model.Image `yaml:"clusterProportionalAutoscalerImage,omitempty"`
        KubeDnsImage                       model.Image `yaml:"kubeDnsImage,omitempty"`
        KubeDnsMasqImage                   model.Image `yaml:"kubeDnsMasqImage,omitempty"`
        KubeReschedulerImage               model.Image `yaml:"kubeReschedulerImage,omitempty"`
        DnsMasqMetricsImage                model.Image `yaml:"dnsMasqMetricsImage,omitempty"`
        ExecHealthzImage                   model.Image `yaml:"execHealthzImage,omitempty"`
        HeapsterImage                      model.Image `yaml:"heapsterImage,omitempty"`
        AddonResizerImage                  model.Image `yaml:"addonResizerImage,omitempty"`
        KubeDashboardImage                 model.Image `yaml:"kubeDashboardImage,omitempty"`
        PauseImage                         model.Image `yaml:"pauseImage,omitempty"`
        FlannelImage                       model.Image `yaml:"flannelImage,omitempty"`
        DexImage                           model.Image `yaml:"dexImage,omitempty"`
        JournaldCloudWatchLogsImage        model.Image `yaml:"journaldCloudWatchLogsImage,omitempty"`
}

type KubernetesContainerImages struct {
	KubernetesVersion string `yaml:"kubernetesVersion"`
	ContainerImages `yaml:",inline"`
}

type WaitSignalConf struct {
	Enabled      bool
	MaxBatchSize uint `yaml:"maxBatchSize"`//validate: >0
}

type NodepoolConf struct {
	InstanceCommonDescrEmbed `yaml:",inline"`
	Name             NodePoolName
	LoadBalancer     struct {
		Enabled bool
		Names   []ec2.ELBName
		//SecurityGroupIds []ec2.SecurityGroup -- removed, duplicate of SGs above
	}
	ApiEndpointName APIEndpointName
	TargetGroup     struct {
		Enabled bool
		Arns    []ec2.ALBTargetGroupARN
		//SecurityGroupIds []ec2.SecurityGroup -- removed, duplicate of SGs above
	}
	VolumeMounts              []VolumeMountConf
	NodeStatusUpdateFrequency kubernetes.TimePeriod // 10s, 5h
	Gpu                       GpuConf

	//SpotPrice uint -- removed, is listed in SpotFleetConf
	WaitSignal       WaitSignalConf
	AutoScalingGroup ASGConf `yaml:"autoScalingGroup"`
	SpotFleet        SpotFleetConf

	Autoscaling              AutoScalingConf
	ExperimentalNodepoolSubsetConf `yaml:",inline"`
	ElasticFileSystemId      ec2.EFSId
	EphemeralImageStorage    struct{ Enabled bool }
	Kube2IamSupport          struct{ Enabled bool }
	KubeletOpts              kubernetes.KubeletOptionsString

	NodeLabels map[kubernetes.LabelName]kubernetes.LabelValue
	Taints     []struct {
		Key    kubernetes.TaintKey
		Value  kubernetes.TaintValue
		Effect kubernetes.TaintEffect
	}
	KubernetesContainerImages `yaml:",inline"`
	SSHAuthorizedKeys  []types.SSHAuthorizedKey
	CustomSettings map[string]interface{} `yaml:"customSettings"`
}

type EtcdConf struct {
	InstanceCommonDescrEmbed `yaml:",inline"`
	DataVolume       MaybeEncryptedOrEphemeralVolume

	Version                types.EtcdVersion
	Snapshot               struct{ Automated bool }
	DisasterRecovery       struct{ Automated bool }
	MemberIdentityProvider types.EtcdMemberIdentityProvider `yaml:"memberIdentityProvider"`
	InternalDomainName     types.DNSName `yaml:"internalDomainName"`
	ManageRecordSets       bool `yaml:"manageRecordSets"`
	HostedZone             struct{ Id ec2.HostedZoneId } `yaml:"hostedZone"`
	Nodes                  []struct {
		Name types.EtcdMemberIdentifier
		Fqdn types.DNSName
	}

	KMSKeyArn          ec2.KMSKeyARN
}

type AutoScalingConf struct {
	ClusterAutoScaler struct {
		Enabled bool
	} `yaml:"clusterAutoScaler"`
}

type EtcAwsEnvironmentConf struct {
	Enabled     bool
	Environment map[string]string
}

type VPCConf struct {
	//VpcId ec2.VPCId -- removed as deprecated
	Vpc struct {
		Id                ec2.VPCId     //conflicts: IdFromStackOutput
		IdFromStackOutput ec2.StackName `yaml:idFromStackOutput` //conflicts: Id;
		//InternetGatewayId ec2.IGWId -- removed as deprecated
	} `yaml:vpc`

	InternetGateway struct {
		Id                ec2.IGWId     //conflicts: IdFromStackOutput
		IdFromStackOutput ec2.StackName `yaml:idFromStackOutput` //conflicts: Id
	}

	// RouteTableId ec2.RouteTableId -- removed in favour of subnets[].routeTable.id

	VpcCIDR types.IPNet `yaml:"vpcCIDR"` //future: should be conflicting with vpcID
	//InstanceCIDR net.IPNet // -- reomved in favour of nodepools[] and subnets[]

	Subnets []SubnetConf
}

type TLSConf struct {
	tlsCADurationDays   uint `default:"3650"`
	tlsCertDurationDays uint `default:"365"`
}

type CloudWatchLoggingConf struct {
	Enabled         bool
	RetentionInDays uint `yaml:"retentionInDays"`
	LocalStreaming  struct {
		Enabled  bool
		Filter   string
		Interval uint
	} `yaml:localStreaming`
}

type AmazonSSMAgentConf struct {
	Enabled     bool
	DownloadURL types.URL
	Sha1Sum     types.SHA1SUM
}

type AuditLogConf struct {
	Enabled bool
	MaxAge  uint
	LogPath string
}

type DexConf struct {
	Enabled         bool
	URL             types.URL
	ClientID        dex.ClientID
	Username        dex.Username
	SelfSignedCa    bool
	Connectors      []map[string]interface{} //TODO better structs
	StaticClients   []map[string]string      //TODO better structs
	StaticPasswords []map[string]string      //TODO better structs
}

type AddonsConf struct {
	ClusterAutoscaler struct{ Enabled bool } `yaml:"clusterAutoscaler"`
	Rescheduler       struct{ Enabled bool } 
}


// subset of experimental flags which can be listed in nodepool config
// (for some reason they are listed inline, not under experimental: field)
type ExperimentalNodepoolSubsetConf struct {
	AwsEnvironment EtcAwsEnvironmentConf
	AwsNodeLabels       struct{ Enabled bool }
	ClusterAutoscalerSupport struct{ Enabled bool } `yaml:"clusterAutoscalerSupport"`
	NodeDrainer         struct{ Enabled bool } `yaml:"nodeDrainer"`
	Kube2IamSupport     struct{ Enabled bool } `yaml:"kube2IamSupport"`
	TLSBootstrap        struct{ Enabled bool } `yaml:"tlsBootstrap"`
}

type ExperimentalConf struct {
	ExperimentalNodepoolSubsetConf `yaml:",inline"`
	Admission struct {
		PodSecurityPolicy  struct{ Enabled bool } `yaml:"podSecurityPolicy"`
		DenyEscalatingExec struct{ Enabled bool } `yaml:"denyEscalatingExec"`
	}

	AuditLog       AuditLogConf
	Authentication struct {
		Webhook struct {
			Enabled      bool
			cacheTTL     kubernetes.TimePeriod
			configBase64 types.Base64Yaml
		}
	}
	EphemeralImageStore struct{ Enabled bool } `yaml:"ephemeralImageStorage"`
	Dex     DexConf
	Plugins struct {
		Rbac struct{ Enabled bool }
	}
	DisableSecurityGroupIngress bool `yaml:"disableSecurityGroupIngress"`
	NodeMonitorGracePeriod      kubernetes.TimePeriod `yaml:"nodeMonitorGracePeriod"`
}

type ClusterYAML struct {
	ClusterName    types.ClusterName `yaml:"clusterName"`
	ReleaseChannel coreos.ReleaseChannel `yaml:"releaseChannel"`

	AmiId                       ec2.AmiId `yaml:"amiId"`
	HostedZoneId                ec2.HostedZoneId `yaml:"hostedZoneId"`
	SshAccessAllowedSourceCIDRs []types.IPNet `yaml:"sshAccessAllowedSourceCIDRs"`

	AdminAPIEndpointName APIEndpointName `yaml:"adminAPIEndpointName"`
	ApiEndpoints         []APIEndpointConf `yaml:"apiEndpoints"`

	KeyName           ec2.SSHKeyPairName `yaml:"keyName"`
	SSHAuthorizedKeys []types.SSHAuthorizedKey `yaml:"sshAuthorizedKeys"`

	Region ec2.Region `yaml:"region"`
	//AvailabilityZone ec2.AvailabilityZone -- removed in favour of nodepools
	KMSKeyArn ec2.KMSKeyARN `yaml:"kmsKeyArn"`

	Controller ControllerConf `yaml:"controller"`

	//WorkerCount uint -- removing
	Worker struct {
		//apiEndpointName  -- removing
		NodePools []NodepoolConf `yaml:"nodePools"`
	} `yaml:"worker"`

	//WorkerCreationTimeout ec2.Timeout  -- removed in favour of nodepools
	//WorkerInstanceType ec2.InstanceType
	//WorkerRootVolumeSize
	//WorkerRootVolumeType
	//WorkerRootVolumeIOPS
	//WorkerTenancy
	//WorkerSpotPrice

	Etcd EtcdConf `yaml:"etcd"`

	VPCConf `yaml:",inline"`

	ServiceCIDR  types.IPNet `yaml:"serviceCIDR"`
	PodCIDR      types.IPNet `yaml:"podCIDR"`
	DnsServiceIP types.IP `yaml:"dnsServiceIP"`
	MapPublicIPs bool  `yaml:"mapPublicIPs"` //future:shouldn't it be per nodepool?

	TLSConf `yaml:",inline"`
	KubernetesContainerImages `yaml:",inline"`

	UseCalico              bool  `yaml:"useCalico"`
	ElasticFileSystemId    ec2.EFSId  `yaml:"elasticFileSystemId"`
	SharedPersistentVolume bool `yaml:"sharedPersistentVolume"`
	ContainerRuntime       types.ContainerRuntime `yaml:"containerRuntime"`

	ManageCertificates      bool  `yaml:"manageCertificates"`
	WaitSignal              WaitSignalConf  `yaml:"waitSignal"` // seems to be related to Controllers only, can be deprecate?
	KubeResourcesAutosave   struct{ Enabled bool } `yaml:"kubeResourcesAutosave"`
	CloudWatchLogging       CloudWatchLoggingConf `yaml:"cloudWatchLogging"`
	AmazonSsmAgent          AmazonSSMAgentConf `yaml:"amazonSsmAgent"`
	KubeDNS                 struct{ NodeLocalResolver bool } `yaml:"kubeDNS"`
	CloudFormationStreaming bool `yaml:"cloudFormationStreaming"`

	Addons AddonsConf
	Experimental ExperimentalConf

	StackTags map[ec2.TagName]ec2.TagValue  `yaml:"stackTags"`

	CustomSettings map[string]interface{}  `yaml:"customSettings"`
}
