package pb

import (
	"fmt"
	"regexp"

	"github.com/fagongzi/gateway/pkg/expr"
	"github.com/fagongzi/gateway/pkg/pb/metapb"
	"github.com/fagongzi/gateway/pkg/plugin"
)

// ValidateRouting validate routing
func ValidateRouting(value *metapb.Routing) error {
	if value.API == 0 {
		return fmt.Errorf("missing api")
	}

	if value.ClusterID == 0 {
		return fmt.Errorf("missing cluster")
	}

	if len(value.Name) == 0 {
		return fmt.Errorf("missing name")
	}

	if value.TrafficRate <= 0 || value.TrafficRate > 100 {
		return fmt.Errorf("error traffic rate: %d", value.TrafficRate)
	}

	return nil
}

// ValidateCluster validate cluster
func ValidateCluster(value *metapb.Cluster) error {
	if len(value.Name) == 0 {
		return fmt.Errorf("missing name")
	}

	return nil
}

// ValidateServer validate server
func ValidateServer(value *metapb.Server) error {
	if len(value.Addr) == 0 {
		return fmt.Errorf("missing server address")
	}

	if value.MaxQPS == 0 {
		return fmt.Errorf("missing server max qps")
	}

	return nil
}

// ValidateAPI validate api
func ValidateAPI(value *metapb.API) error {
	if len(value.Name) == 0 {
		return fmt.Errorf("missing api name")
	}

	if len(value.URLPattern) == 0 {
		return fmt.Errorf("missing URLPattern")
	}

	for _, n := range value.Nodes {
		if len(n.URLRewrite) != 0 {
			_, err := expr.Parse([]byte(n.URLRewrite))
			if err != nil {
				return err
			}
		}

		for _, v := range n.Validations {
			for _, r := range v.Rules {
				if r.RuleType == metapb.RuleRegexp {
					_, err := regexp.Compile(r.Expression)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// ValidatePlugin validate plugin
func ValidatePlugin(value *metapb.Plugin) error {
	if len(value.Name) == 0 {
		return fmt.Errorf("missing plugin name")
	}

	if value.Version == 0 {
		return fmt.Errorf("missing plugin version")
	}

	if len(value.Content) == 0 {
		return fmt.Errorf("missing plugin content")
	}

	_, err := plugin.NewRuntime(value)
	if err != nil {
		return err
	}

	return nil
}
