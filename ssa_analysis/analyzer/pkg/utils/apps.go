package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	APP_PATH_POSTNOTIFICATION = "github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow/postnotification_simple"
	APP_PATH_DIGOTA           = "github.com/blueprint-uservices/blueprint/examples/digota/workflow/digota"
	APP_PATH_DSB_MEDIA_SQL    = "github.com/blueprint-uservices/blueprint/examples/dsb_media_sql/workflow/mediamicroservices_sql"
)

var APPS_SQL_TABLES = map[string][]string{
	// key is the name of the app
	// value is a list of <database_name>:<sql_filepath>
	"dsb_media_sql": {
		"movieid_db:../../blueprint/examples/dsb_media_sql/workflow/mediamicroservices_sql/database/movieid.sql",
		"movieinfo_db:../../blueprint/examples/dsb_media_sql/workflow/mediamicroservices_sql/database/movieinfo.sql",
	},
}

func GetAppDatabaseSQLPaths(app string, autofill bool) (bool, string) {
	if autofill {
		if paths, ok := APPS_SQL_TABLES[app]; ok {
			return true, strings.Join(paths, ";")
		}
		return false, ""
	}

	fmt.Printf("\nPlease specify the sql paths if existent.\nFormat (delimiter is ';'): <database_name>:<sql_path>\n> ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("error reading sql paths for app (%s): %s", app, err.Error())
		return false, ""
	}

	return true, input
}
