package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	specs_app_constraints_referential_integrity "github.com/blueprint-uservices/blueprint/examples/app_constraints_referential_integrity/wiring/specs"
	specs_coupons_app "github.com/blueprint-uservices/blueprint/examples/coupons_app/wiring/specs"
	specs_coupons_app_cache "github.com/blueprint-uservices/blueprint/examples/coupons_app_cache/wiring/specs"
	specs_coupons_app_sql "github.com/blueprint-uservices/blueprint/examples/coupons_app_sql/wiring/specs"
	specs_digota "github.com/blueprint-uservices/blueprint/examples/digota/wiring/specs"
	specs_dsb_hotel "github.com/blueprint-uservices/blueprint/examples/dsb_hotel/wiring/specs"
	specs_dsb_media "github.com/blueprint-uservices/blueprint/examples/dsb_media/wiring/specs"
	specs_dsb_media_sql "github.com/blueprint-uservices/blueprint/examples/dsb_media_sql/wiring/specs"
	specs_dsb_sn "github.com/blueprint-uservices/blueprint/examples/dsb_sn/wiring/specs"
	specs_employee_app "github.com/blueprint-uservices/blueprint/examples/employee_app/wiring/specs"
	specs_foobar "github.com/blueprint-uservices/blueprint/examples/foobar/wiring/specs"
	specs_postnotification "github.com/blueprint-uservices/blueprint/examples/postnotification/wiring/specs"
	specs_postnotification_simple "github.com/blueprint-uservices/blueprint/examples/postnotification_simple/wiring/specs"
	specs_shopping_app "github.com/blueprint-uservices/blueprint/examples/shopping_app/wiring/specs"
	specs_shopping_simple "github.com/blueprint-uservices/blueprint/examples/shopping_simple/wiring/specs"
	specs_sockshop2 "github.com/blueprint-uservices/blueprint/examples/sockshop2/wiring/specs"
	specs_trainticket "github.com/blueprint-uservices/blueprint/examples/train_ticket/wiring/specs"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"

	"analyzer/pkg/logger"
)

const (
	BLUEPRINT_PATH_EXAMPLES     string = "github.com/blueprint-uservices/blueprint/examples/"
	BLUEPRINT_PATH_CORE_BACKEND string = "github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

var APPS_SQL_TABLES = map[string][]string{
	// <database_name>:<sql_filepath>
	"coupons_app_sql": {
		"coupons_db:blueprint/examples/coupons_app_sql/workflow/coupons_app_sql/database/coupons.sql",
		"students_db:blueprint/examples/coupons_app_sql/workflow/coupons_app_sql/database/students.sql",
	},
	"dsb_media_sql": {
		"movieid_db:blueprint/examples/dsb_media_sql/workflow/mediamicroservices_sql/database/movieid.sql",
		"movieinfo_db:blueprint/examples/dsb_media_sql/workflow/mediamicroservices_sql/database/movieinfo.sql",
	},
}

var APPS_MONGODB_SCHEMAS = map[string][]string{
	// <database_name>:<collection_name>:<json_filepath>
	"coupons_app": {
		"coupons_db:Coupon:blueprint/examples/coupons_app/workflow/coupons_app/database/coupons.json",
		"coupons_db:ClaimedCoupon:blueprint/examples/coupons_app/workflow/coupons_app/database/claimed_coupons.json",
		"students_db:Student:blueprint/examples/coupons_app/workflow/coupons_app/database/students.json",
	},
}

var Apps = []string{
	"foobar",
	"shopping_simple",
	"shopping_app",
	"postnotification_simple",
	"postnotification",
	"sockshop2",
	"trainticket",
	"app_constraints_referential_integrity",
	"employee_app",
	"dsb_sn",
	"dsb_hotel",
	"dsb_media",
	"dsb_media_sql",
	"coupons_app",
	"coupons_app_sql",
	"coupons_app_cache",
	"digota",
}

type AppInfo struct {
	PackagePath   string
	BlueprintSpec cmdbuilder.SpecOption
}

var APPS_INFO = map[string]AppInfo{
	"postnotification":                      {BLUEPRINT_PATH_EXAMPLES + "postnotification/workflow/postnotification", specs_postnotification.Docker},
	"postnotification_simple":               {BLUEPRINT_PATH_EXAMPLES + "postnotification_simple/workflow/postnotification_simple", specs_postnotification_simple.Docker},
	"app_constraints_referential_integrity": {BLUEPRINT_PATH_EXAMPLES + "app_constraints_referential_integrity/workflow/app_constraints_referential_integrity", specs_app_constraints_referential_integrity.Docker},
	"employee_app":                          {BLUEPRINT_PATH_EXAMPLES + "employee_app/workflow/employee_app", specs_employee_app.Docker},
	"digota":                                {BLUEPRINT_PATH_EXAMPLES + "digota/workflow/digota", specs_digota.Docker},
	"coupons_app":                           {BLUEPRINT_PATH_EXAMPLES + "coupons_app/workflow/coupons_app", specs_coupons_app.Docker},
	"coupons_app_sql":                       {BLUEPRINT_PATH_EXAMPLES + "coupons_app_sql/workflow/coupons_app_sql", specs_coupons_app_sql.Docker},
	"coupons_app_cache":                     {BLUEPRINT_PATH_EXAMPLES + "coupons_app_cache/workflow/coupons_app_cache", specs_coupons_app_cache.Docker},
	"foobar":                                {BLUEPRINT_PATH_EXAMPLES + "foobar/workflow/foobar", specs_foobar.Docker},
	"sockshop2":                             {BLUEPRINT_PATH_EXAMPLES + "sockshop2/workflow", specs_sockshop2.Docker},
	"trainticket":                           {BLUEPRINT_PATH_EXAMPLES + "train_ticket/workflow", specs_trainticket.Docker},
	"shopping_app":                          {BLUEPRINT_PATH_EXAMPLES + "shopping_app/workflow", specs_shopping_app.Docker},
	"shopping_simple":                       {BLUEPRINT_PATH_EXAMPLES + "shopping_simple/workflow", specs_shopping_simple.Docker},
	"dsb_hotel":                             {BLUEPRINT_PATH_EXAMPLES + "dsb_hotel/workflow/hotelreservation", specs_dsb_hotel.Original},
	"dsb_sn":                                {BLUEPRINT_PATH_EXAMPLES + "dsb_sn/workflow/socialnetwork", specs_dsb_sn.Docker},
	"dsb_media":                             {BLUEPRINT_PATH_EXAMPLES + "dsb_media/workflow/mediamicroservices", specs_dsb_media.Docker},
	"dsb_media_sql":                         {BLUEPRINT_PATH_EXAMPLES + "dsb_media_sql/workflow/mediamicroservices_sql", specs_dsb_media_sql.Docker},
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
		logger.Logger.Fatalf("error reading sql paths for app (%s): %s", app, err.Error())
		return false, ""
	}

	return true, input
}

func GetAppDatabaseDocPaths(app string, autofill bool) (bool, string) {
	if autofill {
		if paths, ok := APPS_MONGODB_SCHEMAS[app]; ok {
			return true, strings.Join(paths, ";")
		}
		return false, ""
	}

	fmt.Printf("\nPlease specify the sql paths if existent.\nFormat (delimiter is ';'): <database_name>:<sql_path>\n> ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		logger.Logger.Fatalf("error reading sql paths for app (%s): %s", app, err.Error())
		return false, ""
	}

	return true, input
}

var APPS_CONSTRAINTS_UNICITY = map[string][]string{
	// <database_name>.<root_object>.<unique_field>
	"coupons_app": {
		"(STUDENTS_DB.Student.StudentID)",
		"(COUPONS_DB.Coupon.CouponID)",
		"(COUPONS_DB.ClaimedCoupon.CouponID,COUPONS_DB.ClaimedCoupon.UserID)",
	},
	"dsb_media": {
		//"(MOVIEID_DB.MovieId.MovieID)",
		"(MOVIEID_DB.MovieId.Title)",
		//"(MOVIEINFO_DB.MovieInfo.MovieID)",
	},
}

func GetAppDatabaseUnicityConstraintFromUserInput(app string, autofill bool) (bool, string) {
	if autofill {
		if constraints, ok := APPS_CONSTRAINTS_UNICITY[app]; ok {
			return true, strings.Join(constraints, ";")
		}
		return false, ""
	}
	fmt.Printf("\nPlease specify fields to enforce unicity constraint.\nFormat (delimiter is ';'): (<unique_field>)[;(<composed_unique_field_1, composed_unique_field_2>)]\n> ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		logger.Logger.Fatalf("error reading sql paths for app (%s): %s", app, err.Error())
		return false, ""
	}
	return true, input
}

func LoadAppPath(app string) string {
	if info, ok := APPS_INFO[app]; ok {
		return info.PackagePath
	}
	logger.Logger.Fatalf("unknown application name %s", app)
	return ""
}

func LoadAppSpec(app string) cmdbuilder.SpecOption {
	if info, ok := APPS_INFO[app]; ok {
		return info.BlueprintSpec
	}
	logger.Logger.Fatalf("unknown application %s", app)
	return cmdbuilder.SpecOption{}
}
