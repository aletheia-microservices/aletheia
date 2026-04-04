package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var registeredAppPaths = map[string]bool{}
var APPS_NOSQL_SCHEMAS = map[string]string{}
var APPS_SQL_TABLES = map[string][]string{}

func RegisterApp(pkgPath string) {
	registeredAppPaths[pkgPath] = true
}

func GetAppRootPackagePath(appsimplename string) string {
	return fmt.Sprintf("github.com/blueprint-uservices/blueprint/examples/%s/workflow/...", appsimplename)
}

func IsAppPackagePath(pkgpath string) bool {
	return registeredAppPaths[pkgpath]
}

func GetAppDatabaseSQLPaths(app string, autofill bool) (bool, string) {
	if autofill {
		if paths, ok := APPS_SQL_TABLES[app]; ok {
			return true, strings.Join(paths, ";")
		}
		return false, ""
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		logrus.Fatalf("error reading sql paths for app (%s): %s", app, err.Error())
		return false, ""
	}

	return true, input
}
