//
// Go OrientDB Driver - An OrientDB-Driver for Go's database/sql package
//
// The driver should be used via the database/sql package:
//
// import "database/sql"
// import _ "github.com/quux00/ogonori/osql"
//
// db, err := sql.Open("ogonori", "admin@admin:127.0.0.1/ogonoriTest")
//
package osql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"

	"github.com/quux00/ogonori/constants"
	"github.com/quux00/ogonori/obinary"
	"github.com/quux00/ogonori/ogl"
)

var dsnRx *regexp.Regexp = regexp.MustCompile(`([^@]+)@([^:]+):([^/]+)/(.+)`)

//
// Implements the Go sql/driver.Driver interface.
//
type OgonoriDriver struct{}

func init() {
	sql.Register("ogonori", &OgonoriDriver{})
}

//
// Open returns a new connection to the database.
// The dsn (driver-specific name) is a string in a driver-specific format.
// For ogonori, the dsn should be of the format:
//   uname@passw:ip-or-host:port/dbname
//   or
//   uname@passw:ip-or-host/dbname  (default port of 2424 is used)
//
func (d *OgonoriDriver) Open(dsn string) (driver.Conn, error) {
	ogl.Debugln("** OgonoriDriver#Open")

	uname, passw, host, port, dbname, err := parseDsn(dsn)
	clientOpts := obinary.ClientOptions{ServerHost: host, ServerPort: port}
	dbc, err := obinary.NewDBClient(clientOpts)
	if err != nil {
		return nil, err
	}

	// TODO: right now assumes DocumentDB type - pass in on the dsn??
	//       NOTE: I tried a graphDB with DocumentDB type and it worked, so why is it necesary at all?
	// TODO: this maybe shouldn't happen in this method -> might do it lazily in Query/Exec methods?
	err = obinary.OpenDatabase(dbc, dbname, constants.DocumentDB, uname, passw)
	if err != nil {
		return nil, err
	}

	return &ogonoriConn{dbc}, nil
}

func parseDsn(dsn string) (uname, passw, host, port, dbname string, err error) {
	matches := dsnRx.FindStringSubmatch(dsn)
	if matches == nil || len(matches) != 5 {
		return "", "", "", "", "",
			fmt.Errorf("Unable to parse connection string: %s. Must be of the format: %s",
				dsn, "uname@passw:ip-or-host/dbname")
	}
	toks := strings.Split(matches[3], ":")
	host = toks[0]
	if len(toks) > 1 {
		port = toks[1]
	}
	return matches[1], matches[2], host, port, matches[4], nil
}
