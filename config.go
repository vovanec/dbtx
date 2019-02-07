package dbtx

type DefaultConfig struct {
	driverName string
	roDSN string
	rwDSN string
	maxOpenConn int
	maxIdleConn int
	errorLog ErrorLogFunc
}

func (c *DefaultConfig) DriverName() string {
	return c.driverName
}

func (c *DefaultConfig) GetRODataSourceName() string {
	return c.roDSN
}

func (c *DefaultConfig) GetRWDataSourceName() string {
	return c.rwDSN
}

func (c *DefaultConfig) MaxOpenConnections() int {
	return c.maxOpenConn
}

func (c *DefaultConfig) MaxIdleConnections() int {
	return c.maxIdleConn
}

func (c *DefaultConfig) ErrorLog() ErrorLogFunc {
	return c.errorLog
}

func (c *DefaultConfig) WithMaxOpenConnections(maxConn int) Config {
	c.maxOpenConn = maxConn
	return c
}

func (c *DefaultConfig) WithMaxIdleConnections(maxConn int) Config {
	c.maxIdleConn = maxConn
	return c
}

func (c *DefaultConfig) WithErrorLog(errLog ErrorLogFunc) Config {
	c.errorLog = errLog
	return c
}

func NewDefaultConfig(driverName, roDSN, rwDSN string) *DefaultConfig {
	return &DefaultConfig{
		driverName: driverName,
		roDSN: roDSN,
		rwDSN: rwDSN,
		errorLog: func (format string, args ...interface{}) {},
	}
}
