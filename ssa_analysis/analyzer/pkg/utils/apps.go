package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

const (
	APP_PATH_POSTNOTIFICATION_SIMPLE = "postnotification_simple/workflow/postnotification_simple"
	APP_PATH_DIGOTA                  = "digota/workflow/digota"
	APP_PATH_DSB_MEDIA_SQL           = "dsb_media_sql/workflow/mediamicroservices_sql"
	APP_PATH_SOCKSHOP3               = "sockshop3/workflow/sockshop3"
	APP_PATH_DSB_SN                  = "dsb_sn/workflow/socialnetwork"
	APP_PATH_FOO_BAR                 = "foobar/workflow/foobar"
	APP_PATH_DSB_HOTEL2              = "dsb_hotel2/workflow/hotelreservation2"
	APP_PATH_TRAIN_TICKET2           = "train_ticket2/workflow/train_ticket2"
)

var APPS_PACKAGE_PATHS = []string{
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_POSTNOTIFICATION_SIMPLE,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_DIGOTA,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_DSB_MEDIA_SQL,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_SOCKSHOP3,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_DSB_SN,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_FOO_BAR,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_DSB_HOTEL2,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_TRAIN_TICKET2,
}

var APPS_SQL_TABLES = map[string][]string{
	// key is the name of the app
	// value is a list of <database_name>:<sql_filepath>
	"dsb_media_sql": {
		"movieid_db:../../blueprint/examples/dsb_media_sql/workflow/mediamicroservices_sql/database/movieid.sql",
		"movieinfo_db:../../blueprint/examples/dsb_media_sql/workflow/mediamicroservices_sql/database/movieinfo.sql",
		"castinfo_db:../../blueprint/examples/dsb_media_sql/workflow/mediamicroservices_sql/database/castinfo.sql",
		"plot_db:../../blueprint/examples/dsb_media_sql/workflow/mediamicroservices_sql/database/plot.sql",
	},
}

func IsAppPackagePath(pkgpath string) bool {
	return slices.Contains(APPS_PACKAGE_PATHS, pkgpath)
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
