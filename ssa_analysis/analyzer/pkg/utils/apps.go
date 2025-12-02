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
	APP_PATH_POSTNOTIFICATION       = "postnotification/workflow/postnotification"
	APP_PATH_DIGOTA                 = "digota/workflow/digota"
	APP_PATH_ESHOPMICROSERVICES     = "eshopmicroservices/workflow/eshopmicroservices"
	APP_PATH_DSB_MEDIAMICROSERVICES = "dsb_mediamicroservices/workflow/mediamicroservices"
	APP_PATH_SOCKSHOP               = "sockshop/workflow/sockshop"
	APP_PATH_DSB_SOCIALNETWORK      = "dsb_socialnetwork/workflow/socialnetwork"
	APP_PATH_FOO_BAR                = "foobar/workflow/foobar"
	APP_PATH_FOO_BAR2               = "foobar2/workflow/foobar2"
	APP_PATH_DSB_HOTEL2             = "dsb_hotel2/workflow/hotelreservation2"
	APP_PATH_TRAIN_TICKET           = "train_ticket/workflow/train_ticket"
	APP_PATH_SYNTHETIC_APP          = "synthetic_app/workflow/synthetic_app"
	APP_PATH_SYNTHETIC_APPA         = "synthetic_app/workflow/synthetic_appA"
	APP_PATH_SYNTHETIC_APPB         = "synthetic_app/workflow/synthetic_appB"
	APP_PATH_SYNTHETIC_APP1         = "synthetic_app1/workflow/synthetic_app1"
	APP_PATH_SYNTHETIC_APP2         = "synthetic_app2/workflow/synthetic_app2"
	APP_PATH_SYNTHETIC_APP3         = "synthetic_app3/workflow/synthetic_app3"
	APP_PATH_SYNTHETIC_APP4         = "synthetic_app4/workflow/synthetic_app4"
	APP_PATH_SYNTHETIC_APP5         = "synthetic_app5/workflow/synthetic_app5"
	APP_PATH_SYNTHETIC_APP6         = "synthetic_app5/workflow/synthetic_app6"
	APP_PATH_SYNTHETIC_APP7         = "synthetic_app5/workflow/synthetic_app7"

	BLUEPRINT_EXAMPLES_RELATIVE_PATH = "../../blueprint/examples"
)

var APPS_PACKAGE_PATHS = []string{
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_POSTNOTIFICATION,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_DIGOTA,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_ESHOPMICROSERVICES,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_DSB_MEDIAMICROSERVICES,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_SOCKSHOP,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_DSB_SOCIALNETWORK,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_FOO_BAR,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_DSB_HOTEL2,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_TRAIN_TICKET,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_SYNTHETIC_APP,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_SYNTHETIC_APPA,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_SYNTHETIC_APPB,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_SYNTHETIC_APP1,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_SYNTHETIC_APP2,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_SYNTHETIC_APP3,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_SYNTHETIC_APP4,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_SYNTHETIC_APP5,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_SYNTHETIC_APP6,
	"github.com/blueprint-uservices/blueprint/examples/" + APP_PATH_SYNTHETIC_APP7,
}

var APPS_NOSQL_SCHEMAS = map[string]string{
	"dsb_mediamicroservices": BLUEPRINT_EXAMPLES_RELATIVE_PATH + "/" + APP_PATH_DSB_MEDIAMICROSERVICES,
}

var APPS_SQL_TABLES = map[string][]string{
	// key is the name of the app
	// value is a list of <database_name>:<sql_filepath>
	"dsb_mediamicroservices_sql": {
		"movieid_db:../../blueprint/examples/dsb_mediamicroservices_sql/workflow/mediamicroservices_sql/database/movieid.sql",
		"movieinfo_db:../../blueprint/examples/dsb_mediamicroservices_sql/workflow/mediamicroservices_sql/database/movieinfo.sql",
		"castinfo_db:../../blueprint/examples/dsb_mediamicroservices_sql/workflow/mediamicroservices_sql/database/castinfo.sql",
		"plot_db:../../blueprint/examples/dsb_mediamicroservices_sql/workflow/mediamicroservices_sql/database/plot.sql",
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

	// EVAL: logrus.Tracef("\nPlease specify the sql paths if existent.\nFormat (delimiter is ';'): <database_name>:<sql_path>\n> ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		logrus.Fatalf("error reading sql paths for app (%s): %s", app, err.Error())
		return false, ""
	}

	return true, input
}
