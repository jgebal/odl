package db

import (
	"net/http"

	"github.com/tschf/odl/types"
	"github.com/tschf/odl/types/arch"
)

func GetXeResouces() []*types.Resource {

	acceptCookie := &http.Cookie{
		Name:   "oraclelicense",
		Value:  "accept-sqldev-cookie",
		Domain: ".oracle.com",
	}

	xeResources := []*types.Resource{}

	// Oracle 11gXE for Linux, 64-bit
	xeResources = append(xeResources, &types.Resource{
		Component:    "db",
		Version:      "11gXE",
		File:         "https://edelivery.oracle.com/akam/otn/linux/oracle11g/xe/oracle-xe-11.2.0-1.0.x86_64.rpm.zip",
		License:      "http://www.oracle.com/technetwork/licenses/database-11g-express-license-459621.html",
		OS:           "linux",
		Arch:         arch.X64,
		Lang:         "na",
		AcceptCookie: acceptCookie,
	})

	// Oracle 11gXE for Windows, 32-bit
	xeResources = append(xeResources, &types.Resource{
		Component:    "db",
		Version:      "11gXE",
		File:         "https://edelivery.oracle.com/akam/otn/nt/oracle11g/xe/OracleXE112_Win32.zip",
		License:      "http://www.oracle.com/technetwork/licenses/database-11g-express-license-459621.html",
		OS:           "windows",
		Arch:         arch.X86,
		Lang:         "na",
		AcceptCookie: acceptCookie,
	})

	// Oracle 11gXE for windows, 64bit
	xeResources = append(xeResources, &types.Resource{
		Component:    "db",
		Version:      "11gXE",
		File:         "https://edelivery.oracle.com/akam/otn/nt/oracle11g/xe/OracleXE112_Win64.zip",
		License:      "http://www.oracle.com/technetwork/licenses/database-11g-express-license-459621.html",
		OS:           "windows",
		Arch:         arch.X64,
		Lang:         "na",
		AcceptCookie: acceptCookie,
	})

	return xeResources
}
