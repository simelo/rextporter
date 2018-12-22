package tomlconfig

import (
	"errors"

	"github.com/simelo/rextporter/src/configlocator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// ErrKeyEmptyValue check for empty and/or null
	ErrKeyEmptyValue = errors.New("A required value is missed")
	// ErrKeyReadingFile check for empty and/or null
	ErrKeyReadingFile = errors.New("Error reading file")
	// ErrKeyDecodingFile check for empty and/or null
	ErrKeyDecodingFile = errors.New("Error decoding read content")
)

type configFromFile struct {
	filePath string
}

type mainConfig struct {
	ServicesConfigTransport          string
	ServicesConfigPath               string
	MetricsForServicesConfigPath     string
	ResourcePathsForServicesConfPath string
}

func (cf configFromFile) readMainConf() (mainConf mainConfig, err error) {
	if len(cf.filePath) == 0 {
		log.WithError(ErrKeyEmptyValue).Errorln("file path is required to read main config")
		return mainConf, ErrKeyEmptyValue
	}
	viper.SetConfigFile(cf.filePath)
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Errorln("error reading main toml config")
		return mainConf, ErrKeyReadingFile
	}
	if err = viper.Unmarshal(&mainConf); err != nil {
		log.WithFields(log.Fields{"err": err, "path": cf.filePath}).Errorln("Error decoding main config file content")
		return mainConf, ErrKeyReadingFile
	}
	return mainConf, nil
}

func (cf configFromFile) readServicesConf() (services []Service, err error) {
	if len(cf.filePath) == 0 {
		log.WithError(ErrKeyEmptyValue).Errorln("file path is required to read services config")
		return services, ErrKeyEmptyValue
	}
	viper.SetConfigFile(cf.filePath)
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Errorln("error reading toml services toml config")
		return services, ErrKeyReadingFile
	}
	var root RootConfig
	if err := viper.Unmarshal(&root); err != nil {
		log.WithFields(log.Fields{"err": err, "path": cf.filePath}).Errorln("Error decoding services config file content")
		return services, ErrKeyReadingFile
	}
	services = root.Services
	return services, err
}

func (cf configFromFile) readMetricsForServiceConf() (metricsConf MetricsTemplate, err error) {
	if len(cf.filePath) == 0 {
		log.WithError(ErrKeyEmptyValue).Errorln("file path is required to read metrics for service config")
		return metricsConf, ErrKeyEmptyValue
	}
	viper.SetConfigFile(cf.filePath)
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Errorln("error reading metrics for service toml config")
		return metricsConf, ErrKeyReadingFile
	}
	type metricsForServiceConfig struct {
		Metrics MetricsTemplate
	}
	var metricsForServiceConf metricsForServiceConfig
	if err := viper.Unmarshal(&metricsForServiceConf); err != nil {
		log.WithFields(log.Fields{"err": err, "path": cf.filePath}).Errorln("Error decoding metrics for service config file content")
		return metricsConf, ErrKeyReadingFile
	}
	metricsConf = metricsForServiceConf.Metrics
	return metricsConf, err
}

func (cf configFromFile) readResourcePathsForServiceConf() (resPaths4Service ResourcePathTemplate, err error) {
	if len(cf.filePath) == 0 {
		log.WithError(ErrKeyEmptyValue).Errorln("file path is required to read resources paths for service config")
		return resPaths4Service, ErrKeyEmptyValue
	}
	viper.SetConfigFile(cf.filePath)
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Errorln("error reading resources paths for service toml config")
		return resPaths4Service, ErrKeyReadingFile
	}
	type resourcePathsForServiceConfig struct {
		ResourcePaths ResourcePathTemplate
	}
	var resourcePathsForServiceConf resourcePathsForServiceConfig
	if err := viper.Unmarshal(&resourcePathsForServiceConf); err != nil {
		log.WithFields(log.Fields{"err": err, "path": cf.filePath}).Errorln("Error decoding resources paths for service config file content")
		return resPaths4Service, ErrKeyReadingFile
	}
	resPaths4Service = resourcePathsForServiceConf.ResourcePaths
	return resPaths4Service, err
}

func (cf configFromFile) readResourcePathsForServicesConf() (resPaths4Services map[string]string, err error) {
	if len(cf.filePath) == 0 {
		log.WithError(ErrKeyEmptyValue).Errorln("file path is required to read resources paths for services config")
		return resPaths4Services, ErrKeyEmptyValue
	}
	viper.SetConfigFile(cf.filePath)
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Errorln("error reading resources paths for services toml config")
		return resPaths4Services, ErrKeyReadingFile
	}
	type resourcePathsForServicesConfig struct {
		ResourcePathsForServicesConfig map[string]string
	}
	var resourcePathsForServicesConf resourcePathsForServicesConfig
	if err := viper.Unmarshal(&resourcePathsForServicesConf); err != nil {
		log.WithFields(log.Fields{"err": err, "path": cf.filePath}).Errorln("Error decoding resources paths for services config file content")
		return resPaths4Services, ErrKeyReadingFile
	}
	resPaths4Services = resourcePathsForServicesConf.ResourcePathsForServicesConfig
	return resPaths4Services, err
}

func (cf configFromFile) readMetricsPathsForServicesConf() (mtrPaths4Services map[string]string, err error) {
	if len(cf.filePath) == 0 {
		log.WithError(ErrKeyEmptyValue).Errorln("file path is required to read metric paths for services config")
		return mtrPaths4Services, ErrKeyEmptyValue
	}
	viper.SetConfigFile(cf.filePath)
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Errorln("error reading resources paths for services toml config")
		return mtrPaths4Services, ErrKeyReadingFile
	}
	type metricPathsForServicesConfig struct {
		MetricPathsForServicesConfig map[string]string
	}
	var metricPathsForServicesConf metricPathsForServicesConfig
	if err := viper.Unmarshal(&metricPathsForServicesConf); err != nil {
		log.WithFields(log.Fields{"err": err, "path": cf.filePath}).Errorln("Error decoding resources paths for services config file content")
		return mtrPaths4Services, ErrKeyReadingFile
	}
	mtrPaths4Services = metricPathsForServicesConf.MetricPathsForServicesConfig
	return mtrPaths4Services, err
}

func readRootStructure(mainConf mainConfig) (rootConf RootConfig, err error) {
	srvConfReader := configFromFile{filePath: mainConf.ServicesConfigPath}
	rootConf.Services, err = srvConfReader.readServicesConf()
	if err != nil {
		log.WithError(err).Errorln("error reading services config")
		return rootConf, ErrKeyReadingFile
	}
	resPaths4ServicesReader := configFromFile{filePath: mainConf.ResourcePathsForServicesConfPath}
	var resPath4Service, mtrPath4Service map[string]string
	if resPath4Service, err = resPaths4ServicesReader.readResourcePathsForServicesConf(); err != nil {
		log.WithError(err).Errorln("error reading resource paths for services config")
		return rootConf, ErrKeyReadingFile
	}
	mtrPaths4ServicesReader := configFromFile{filePath: mainConf.MetricsForServicesConfigPath}
	if mtrPath4Service, err = mtrPaths4ServicesReader.readMetricsPathsForServicesConf(); err != nil {
		log.WithError(err).Errorln("error reading metric paths for services config")
		return rootConf, ErrKeyReadingFile
	}
	for idxService, service := range rootConf.Services {
		resPath4ServiceReader := configFromFile{filePath: resPath4Service[service.Name]}
		if rootConf.Services[idxService].ResourcePaths, err = resPath4ServiceReader.readResourcePathsForServiceConf(); err != nil {
			log.WithFields(log.Fields{"err": err, "service": service.Name}).Warnln("error reading resource paths for service")
			return rootConf, ErrKeyReadingFile
		}
		mtrPath4ServiceReader := configFromFile{filePath: mtrPath4Service[service.Name]}
		if rootConf.Services[idxService].Metrics, err = mtrPath4ServiceReader.readMetricsForServiceConf(); err != nil {
			log.WithFields(log.Fields{"err": err, "service": service.Name}).Warnln("error reading resource paths for service")
			return rootConf, ErrKeyReadingFile
		}
	}
	return rootConf, err
}

// ReadConfigFromFileSystem will read the config from the file system.
func ReadConfigFromFileSystem(filePath string) (rootConf RootConfig, err error) {
	const generalScopeErr = "error getting config values from file system"
	mainConfigPath := configlocator.MainFile()
	if len(filePath) != 0 {
		mainConfigPath = filePath
	}
	mainConfReader := configFromFile{filePath: mainConfigPath}
	var mainConf mainConfig
	if mainConf, err = mainConfReader.readMainConf(); err != nil {
		log.WithError(err).Errorln("error reading main config file")
		return rootConf, ErrKeyReadingFile
	}
	rootConf, err = readRootStructure(mainConf)
	if err != nil {
		log.WithError(err).Errorln("error reading root structure conf")
		return rootConf, err
	}
	return rootConf, nil
}
