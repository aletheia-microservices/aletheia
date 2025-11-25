package utils

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	APP_PATH_POSTNOTIFICATION_SIMPLE = "postnotification_simple/workflow/postnotification_simple"
	APP_PATH_DIGOTA                  = "digota/workflow/digota"
	APP_PATH_ESHOPMICROSERVICES      = "eshopmicroservices/workflow/eshopmicroservices"
	APP_PATH_DSB_MEDIA_SQL           = "dsb_media_sql/workflow/mediamicroservices_sql"
	APP_PATH_DSB_MEDIA_NOSQL         = "dsb_media_nosql/workflow/mediamicroservices_nosql"
	APP_PATH_SOCKSHOP3               = "sockshop3/workflow/sockshop3"
	APP_PATH_DSB_SN2                 = "dsb_sn2/workflow/socialnetwork2"
	APP_PATH_FOO_BAR                 = "foobar/workflow/foobar"
	APP_PATH_FOO_BAR2                = "foobar2/workflow/foobar2"
	APP_PATH_DSB_HOTEL2              = "dsb_hotel2/workflow/hotelreservation2"
	APP_PATH_TRAIN_TICKET2           = "train_ticket2/workflow/train_ticket2"
	APP_PATH_LARGE_SCALE_APP         = "large_scale_app/workflow/large_scale_app"
	APP_PATH_LARGE_SCALE_APP_A       = "large_scale_app_A/workflow/large_scale_app_A"
	APP_PATH_LARGE_SCALE_APP_B       = "large_scale_app_B/workflow/large_scale_app_B"
	APP_PATH_LARGE_SCALE_APP_C       = "large_scale_app_C/workflow/large_scale_app_C"
	APP_PATH_LARGE_SCALE_APP_D       = "large_scale_app_D/workflow/large_scale_app_D"
	APP_PATH_LARGE_SCALE_APP_E       = "large_scale_app_E/workflow/large_scale_app_E"

	BLUEPRINT_EXAMPLES_RELATIVE_PATH = "../../blueprint/examples"
)

var APPS_PACKAGE_PATHS = []string{
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_POSTNOTIFICATION_SIMPLE,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_DIGOTA,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_ESHOPMICROSERVICES,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_DSB_MEDIA_SQL,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_DSB_MEDIA_NOSQL,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_SOCKSHOP3,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_DSB_SN2,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_FOO_BAR,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_DSB_HOTEL2,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_TRAIN_TICKET2,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_LARGE_SCALE_APP,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_LARGE_SCALE_APP_A,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_LARGE_SCALE_APP_B,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_LARGE_SCALE_APP_C,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_LARGE_SCALE_APP_D,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_LARGE_SCALE_APP_E,
}

var APPS_NOSQL_SCHEMAS = map[string]string{
	"dsb_media_nosql": BLUEPRINT_EXAMPLES_RELATIVE_PATH + "/" + APP_PATH_DSB_MEDIA_NOSQL,
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

func GetAppRootPackagePath(appsimplename string) string {
	return fmt.Sprintf("github.com/blueprint-uservices/blueprint/examples/%s/workflow/...", appsimplename)
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

	// EVAL: fmt.Printf("\nPlease specify the sql paths if existent.\nFormat (delimiter is ';'): <database_name>:<sql_path>\n> ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		logrus.Fatalf("error reading sql paths for app (%s): %s", app, err.Error())
		return false, ""
	}

	return true, input
}
