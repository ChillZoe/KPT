package observe

import (
	"strings"

	"github.com/KPT/conf"
	"github.com/idoubi/goz"
	log "github.com/sirupsen/logrus"
)

func CheckCloudMetadataAPI() {
	for _, apiInstance := range conf.CloudAPI {
		cli := goz.NewClient(goz.Options{
			Timeout: 1,
		})
		resp, err := cli.Get(apiInstance.API)
		if err != nil {
			log.WithFields(log.Fields{"SUCCESS": false}).Warn("Not find %s Metadata API!", apiInstance.CloudProvider)
			continue
		}
		r, _ := resp.GetBody()
		if strings.Contains(r.String(), apiInstance.ResponseMatch) {
			log.WithFields(log.Fields{"SUCCESS": true}).Info("\t%s Metadata API available in %s\n", apiInstance.CloudProvider, apiInstance.API)
			log.Info("\tDocs: %s\n", apiInstance.DocURL)
		} else {
			log.WithFields(log.Fields{"SUCCESS": false}).Warn("Not find %s API!", apiInstance.CloudProvider)
		}
	}
}
