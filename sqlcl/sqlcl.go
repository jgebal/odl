package sqlcl

import (
	"net/http"

	"github.com/tschf/odl/types"
	"github.com/tschf/odl/types/arch"
)

func GetSqlclResources() []*types.Resource {

	acceptCookie := &http.Cookie{
		Name:   "oraclelicense",
		Value:  "accept-sqldev-cookie",
		Domain: ".oracle.com",
	}

	sqlClResources := []*types.Resource{}

	sqlClResources = append(sqlClResources, &types.Resource{
		Component:    "sqlcl",
		Version:      "4.2",
		File:         "https://edelivery.oracle.com/akam/otn/java/sqldeveloper/sqlcl-4.2.0.16.355.0402-no-jre.zip",
		License:      "http://www.oracle.com/technetwork/licenses/sqldev-license-152021.html",
		OS:           "na",
		Arch:         arch.Na,
		Lang:         "na",
		AcceptCookie: acceptCookie,
	})

	return sqlClResources
}
