package proxy

import (
	"errors"
	"plugin"
	"strings"

	"github.com/fagongzi/gateway/pkg/filter"
)

var (
	// ErrUnknownFilter unknown filter error
	ErrUnknownFilter = errors.New("unknown filter")
)

const (
	// FilterBase base filter
	FilterBase = "Base"
	// FilterPrepare prepare filter
	FilterPrepare = "PREPARE"
	// FilterHTTPAccess access log filter
	FilterHTTPAccess = "HTTP-ACCESS"
	// FilterHeader header filter
	FilterHeader = "HEADER" // process header filter
	// FilterXForward x-forward filter
	FilterXForward = "X-FORWARD"
	// FilterBlackList blacklist filter
	FilterBlackList = "BLACKLIST"
	// FilterWhiteList whitelist filter
	FilterWhiteList = "WHITELIST"
	// FilterAnalysis analysis filter
	FilterAnalysis = "ANALYSIS"
	// FilterRateLimiting limit filter
	FilterRateLimiting = "RATE-LIMITING"
	// FilterCircuitBreaker circuit breaker filter
	FilterCircuitBreaker = "CIRCUIT-BREAKER"
	// FilterValidation validation request filter
	FilterValidation = "VALIDATION"
	// FilterCaching caching filter
	FilterCaching = "CACHING"
	// FilterJWT jwt filter
	FilterJWT = "JWT"
	// FilterCross cross filter
	FilterCross = "CROSS"
	// FilterJSPlugin js plugin engine
	FilterJSPlugin = "JS-ENGINE"
)

func (p *Proxy) newFilter(filterSpec *FilterSpec) (filter.Filter, error) {
	if filterSpec.External {
		return newExternalFilter(filterSpec)
	}

	input := strings.ToUpper(filterSpec.Name)

	switch input {
	case FilterPrepare:
		return newPrepareFilter(), nil
	case FilterHTTPAccess:
		return newAccessFilter(), nil
	case FilterHeader:
		return newHeadersFilter(), nil
	case FilterXForward:
		return newXForwardForFilter(), nil
	case FilterAnalysis:
		return newAnalysisFilter(), nil
	case FilterBlackList:
		return newBlackListFilter(), nil
	case FilterWhiteList:
		return newWhiteListFilter(), nil
	case FilterRateLimiting:
		return newRateLimitingFilter(), nil
	case FilterCircuitBreaker:
		return newCircuitBreakeFilter(), nil
	case FilterValidation:
		return newValidationFilter(), nil
	case FilterCaching:
		return newCachingFilter(p.cfg.Option.LimitBytesCaching, p.dispatcher.tw), nil
	case FilterJWT:
		return newJWTFilter(p.cfg.Option.JWTCfgFile)
	case FilterCross:
		return newCrossDomainFilter(p.cfg.Option.CrossCfgFile)
	case FilterJSPlugin:
		return p.jsEngine, nil
	default:
		return nil, ErrUnknownFilter
	}
}

func newExternalFilter(filterSpec *FilterSpec) (filter.Filter, error) {
	p, err := plugin.Open(filterSpec.ExternalPluginFile)
	if err != nil {
		return nil, err
	}

	s, err := p.Lookup("NewExternalFilter")
	if err != nil {
		return nil, err
	}

	sf := s.(func() (filter.Filter, error))
	return sf()
}
