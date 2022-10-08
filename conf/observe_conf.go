package conf

var K8sSATokenDefaultPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"

var ScanDefaultPorts = "443,2379,6443,10250,10255"

// Check cloud provider APIs in evaluate task
type cloudAPIS struct {
	CloudProvider string
	API           string
	ResponseMatch string
	DocURL        string
}

var CloudAPI = []cloudAPIS{
	{
		CloudProvider: "Alibaba Cloud",
		API:           "http://100.100.100.200/latest/meta-data/",
		ResponseMatch: "instance-id",
		DocURL:        "https://help.aliyun.com/knowledge_detail/49122.html",
	},
	{
		CloudProvider: "Azure",
		API:           "http://169.254.169.254/metadata/instance",
		ResponseMatch: "azEnvironment",
		DocURL:        "https://docs.microsoft.com/en-us/azure/virtual-machines/windows/instance-metadata-service",
	},
	{
		CloudProvider: "Google Cloud",
		API:           "http://metadata.google.internal/computeMetadata/v1/instance/disks/?recursive=true",
		ResponseMatch: "deviceName",
		DocURL:        "https://cloud.google.com/compute/docs/storing-retrieving-metadata",
	},
	{
		CloudProvider: "Tencent Cloud",
		API:           "http://metadata.tencentyun.com/latest/meta-data/",
		ResponseMatch: "instance-name",
		DocURL:        "https://cloud.tencent.com/document/product/213/4934",
	},
	{
		CloudProvider: "OpenStack",
		API:           "http://169.254.169.254/openstack/latest/meta_data.json",
		ResponseMatch: "availability_zone",
		DocURL:        "https://docs.openstack.org/nova/rocky/user/metadata-service.html",
	},
	{
		CloudProvider: "Amazon Web Services (AWS)",
		API:           "http://169.254.169.254/latest/meta-data/",
		ResponseMatch: "instance-id",
		DocURL:        "https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-retrieval.html",
	},
	{
		CloudProvider: "ucloud",
		API:           "http://100.80.80.80/meta-data/latest/uhost/",
		ResponseMatch: "uhost-id",
		DocURL:        "https://docs.ucloud.cn/uhost/guide/metadata/metadata-server",
	},
}
