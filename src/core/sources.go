package core

import (
	"errors"
)

var (
	// ErrInvalidType for unecpected type
	ErrInvalidType = errors.New("Unsupported type")
	// ErrKeyNotFound in key value store
	ErrKeyNotFound = errors.New("Missing key")
	// ErrNotClonable in key value store
	ErrNotClonable = errors.New("Impossible to obtain a copy of object")
)

const (
	// KeyTypeCounter is the key you should define in the config file for counters.
	KeyTypeCounter = "Counter"
	// KeyTypeGauge is the key you should define in the config file for gauges.
	KeyTypeGauge = "Gauge"
	// KeyTypeHistogram is the key you should define in the config file for histograms.
	KeyTypeHistogram = "Histogram"
	// KeyTypeSummary is the key you should define in the config file for summaries.
	KeyTypeSummary = "Summary"
)

// RextServiceDef encapsulates all data for services
type RextServiceDef interface {
	SetBasePath(path string) // can be an http server base path, a filesystem directory ...
	// file | http | ftp
	GetProtocol() string
	SetProtocol(string)
	SetAuthForBaseURL(RextAuthDef)
	AddSource(source RextResourceDef)
	AddSources(sources ...RextResourceDef)
	GetOptions() RextKeyValueStore
}

// RextDataSourceDef for retrieving raw data
type RextResourceDef interface {
	GetResourcePATH(basePath string) string         // url = base path + uri
	GetAuth(defAuth RextAuthDef) (auth RextAuthDef) // base url if none specific for uri
	// http -> /api/v1/network/connections | /api/v1/health
	// file -> /path/to/a/file, can be a json file, a xml or .rar
	SetResourceURI(string) // can be a filesystem path
	SetAuth(RextAuthDef)
	GetDecoder() RextDecoderDef
	SetDecoder(RextDecoderDef)
	AddRextDataPath(RextDataPathDef)
	GetOptions() RextKeyValueStore
}

type RextDecoderDef interface {
	// json, xml, ini, plain_text, fordwader, .rar(even encrypted)
	GetType() string
	GetOptions() RextKeyValueStore
}

// RextDataSourceDef for retrieving raw data
type RextDataPathDef interface {
	AddMetricDef(RextMetricDef)
	GetNodeSolver() RextNodeSolver
	SetNodeSolver(RextNodeSolver)
	GetOptions() RextKeyValueStore
}

const (
	RextNodeSolverTypeJsonPath = "jsonPath"
)

type RextNodeSolver interface {
	// jpath, xpath, ini, plain_text, .rar
	GetType() string
	// "json" -> "/blockchain/head/seq" | "/blockchain/head/fee"
	// "xml" -> "/blockchain/head/seq" | "/blockchain/head/fee"
	// ini -> "key_name"
	// plain_text -> line number
	GetNodePath() string
	SetNodePath(string)
	GetOptions() RextKeyValueStore
}

// RextMetricDef contains the metadata associated to performance metrics
type RextMetricDef interface {
	GetMetricName() string
	GetMetricType() string
	GetMetricDescription() string
	GetMetricLabels() []RextKeyValueStore
	SetMetricName(string)
	SetMetricType(string)
	SetMetricDescription(string)
	SetMetricLabels([]RextKeyValueStore)
	GetOptions() RextKeyValueStore
}

// RextAuth implements an authentication strategies
type RextAuthDef interface {
	GetAuthType() string
	GetOptions() RextKeyValueStore
}

// RextKeyValueStore providing access to object settings
type RextKeyValueStore interface {
	GetString(key string) (string, error)
	SetString(key string, value string) (bool, error)
	GetObject(key string) (interface{}, error)
	SetObject(key string, value interface{}) (bool, error)
	GetKeys() []string
	Clone() (RextKeyValueStore, error)
}
