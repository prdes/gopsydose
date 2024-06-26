package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/psybits/gopsydose"
)

var (
	drugname = flag.String(
		"drug",
		"none",
		"The name of the drug.\nTry to be as accurate as possible!\n"+
			"This has to match the source information.")

	drugroute = flag.String(
		"route",
		"none",
		"oral, smoked, sublingual, insufflation, inhalation, intravenous, etc.\n"+
			"This has to match the source information.")

	drugargdose = flag.Float64(
		"dose",
		0,
		"Just a number, without any units such as ml around it.")

	drugunits = flag.String(
		"units",
		"none",
		"the units themselves: ml, L, mg etc.\n"+
			"This has to match the source information.")

	drugperc = flag.Float64(
		"perc",
		0,
		"This is used for certain drugs which support unit conversion.\n"+
			"It depends on the source. Again just a number, no % around it.")

	drugcost = flag.Float64(
		"cost",
		0,
		"The cost in money for the logged dose.\n"+
			"This is just a number, the currency is set by -cost-cur\n"+
			"and in the settings file. Currently the cost must be\n"+
			"calculated manually for every dose.")

	costCur = flag.String(
		"cost-cur",
		"",
		"The currency to be used for when logging a cost.\n"+
			"This takes priority over the value set in the settings file.")

	changeLog = flag.Bool(
		"change-log",
		false,
		"Make changes to an entry/log/dosage.\n"+
			"If -for-id is not specified, it's to the newest entry.\n"+
			"Must be used in combination with:\n"+
			"-end-time ; -start-time ; -drug and etc.\n"+
			"in order to clarify what you want to change in the entry.")

	endTime = flag.String(
		"end-time",
		"none",
		"Change the end time of the last log.\n"+
			"It accepts unix timestamps as input.\n"+
			"If input is the string \"now\", it will use the current time\n"+
			"\n"+
			"Must be used in combination with -change-log.\n"+
			"If it's also combined with -for-id, it will change for a specific ID.")

	startTime = flag.String(
		"start-time",
		"none",
		"Change the start time of the last log.\n"+
			"It accepts unix timestamps as input.\n"+
			"If input is the string \"now\", it will use the current time\n"+
			"\n"+
			"Must be used in combination with -change-log.\n"+
			"If it's also combined with -for-id, it will change for a specific ID.")

	forUser = flag.String(
		"user",
		drugdose.DefaultUsername,
		"Log for a specific user.\n"+
			"This takes priority over the value set in the settings file.")

	getNewLogs = flag.Int(
		"get-new-logs",
		0,
		"Print a given number of the newest logs for the set user.")

	getOldLogs = flag.Int(
		"get-old-logs",
		0,
		"Print a given number of the oldest logs for the set user.")

	getLogs = flag.Bool(
		"get-logs",
		false,
		"Print all logs for the set user.")

	getLogsCount = flag.Bool(
		"get-logs-count",
		false,
		"Print the total number of logs for the set user.")

	removeNew = flag.Int(
		"clean-new-logs",
		0,
		"Clean a given number of the newest logs.")

	removeOld = flag.Int(
		"clean-old-logs",
		0,
		"Clean a given number of oldest logs.")

	cleanDB = flag.Bool(
		"clean-db",
		false,
		"Remove all tables from the database.\n"+
			"Remember that it's for the currently set database path and driver.")

	cleanLogs = flag.Bool(
		"clean-logs",
		false,
		"Clean all of the logs for the currently set user.")

	cleanNames = flag.Bool(
		"clean-names",
		false,
		"Clean the alternative names from the database.\n"+
			"This includes the replace names for the\n"+
			"currently configured source.\n"+
			"When you reuse the name matching, it will\n"+
			"recreate the tables with the present config files.")

	overwriteNames = flag.Bool(
		"overwrite-names",
		false,
		"Overwrite the alternative names in the database.\n"+
			"It will delete the old tables and replace them\n"+
			"with the currently present ones in the config directory.")

	cleanInfo = flag.Bool(
		"clean-info",
		false,
		"Clean the currently configured information (source) table.\n"+
			"Meaning all remotely fetched dosage ranges and routes\n"+
			"for all drugs. Keep in mind that if you have configured a\n"+
			"different source earlier, it will not be cleaned, unless\n"+
			"you change the configuration back and use this flag again.")

	forID = flag.Int64(
		"for-id",
		0,
		"Perform an action for a particular ID (start of dose timestamp).")

	sourcecfg = flag.String(
		"sourcecfg",
		drugdose.DefaultSource,
		"The name of the source that you want to initialise\n"+
			"settings and sources config files for.")

	sourceAddress = flag.String(
		"source-address",
		drugdose.DefaultSourceAddress,
		"The address of the source that you want to initialise\n"+
			"sources config file for, combined with -sourcecfg.")

	recreateSettings = flag.Bool(
		"recreate-settings",
		false,
		"Recreate the global settings config file.")

	recreateSources = flag.Bool(
		"recreate-sources",
		false,
		"Recreate the global sources config file.")

	dbDir = flag.String(
		"db-dir",
		drugdose.DefaultDBDir,
		"Full path of the database directory.\n"+
			"This will work only on the initial run.\n"+
			"It can be changed later in the settings config file.\n"+
			"If changed after initial run, the old database directory must be removed.")

	getPaths = flag.Bool(
		"get-paths",
		false,
		"Prints the settings and the database paths.")

	getLocalInfoDrug = flag.String(
		"get-local-info-drug",
		"none",
		"Print all of the information about a given drug name from,\n"+
			"the locally stored information table in the database.\n"+
			"For example if you've forgotten dose ranges, times and etc.")

	getLocalInfoDrugs = flag.Bool(
		"get-local-info-drugs",
		false,
		"Print all unique drug names from the locally stored information table.\n"+
			"The used table is according to the set source.")

	getLoggedDrugs = flag.Bool(
		"get-logged-drugs",
		false,
		"Print all unique drug names from the user logs table,\n"+
			"using the provided username")

	getTotalCosts = flag.Bool(
		"get-total-costs",
		false,
		"Print all costs for all drugs in all currencies from the user logs table.\n"+
			"The information is only for the currently set user.")

	getSubNames = flag.String(
		"get-subst-alt-names",
		"",
		"Get all alternative names for a substance.")

	getRouteNames = flag.String(
		"get-route-alt-names",
		"",
		"Get all alternative names for a route.")

	getUnitsNames = flag.String(
		"get-units-alt-names",
		"",
		"Get all alternative names for a unit.")

	dontLog = flag.Bool(
		"dont-log",
		false,
		"Only fetch info about drug to the local information table,\n"+
			"but don't add in user logs table.")

	removeInfoDrug = flag.String(
		"remove-info-drug",
		"none",
		"Remove all entries of a single drug from the local information table.")

	getTimes = flag.Bool(
		"get-times",
		false,
		"Get the times till onset, comeup, etc. according to the current time.\n"+
			"Can be combined with -for-id to get times for a specific ID,\n"+
			"relative to the current time.")

	getUsers = flag.Bool(
		"get-users",
		false,
		"Get all usernames logged.")

	getDBSize = flag.Bool(
		"get-db-size",
		false,
		"Get total size in MiB and bytes for the currently configured database.")

	stopOnCfgInit = flag.Bool(
		"stop-on-config-init",
		false,
		"Stops the program once the config files have been initialised.")

	stopOnDbInit = flag.Bool(
		"stop-on-db-init",
		false,
		"Stops the program once the database file has been created\n"+
			"and initialised, if it doesn't exists already.")

	verbose = flag.Bool(
		"verbose",
		drugdose.DefaultVerbose,
		"Print extra information.")

	remember = flag.Bool(
		"remember",
		false,
		"Remember the -drug -units and -route.\n"+
			"On the next run, only -dose can be used and it will\n"+
			"reuse the information from the previous run.")

	forget = flag.Bool(
		"forget",
		false,
		"Forget the remembered -drug -units and -route.")

	getDBDriver = flag.Bool(
		"get-db-driver",
		false,
		"Show info about the currently configured database driver.")

	noGetLimit = flag.Bool(
		"no-get-limit",
		false,
		"There is a default limit of getting\n"+
			"a maximum of 100 entries. You can bypass it\n"+
			"using this option. This does not\n"+
			"affect logging new entries, for that do -get-dirs\n"+
			"and checkout the settings file.")

	searchStr = flag.String(
		"search",
		"none",
		"Search all columns for a specific string.")

	searchExact = flag.Bool(
		"search-exact",
		false,
		"Search a specific column.\n"+
			"The column you search is dependent on\n"+
			"the -drug -route -units flags. You don't need\n"+
			"to use the -search flag.\n"+
			"Compared to -search, this doesn't look if the\n"+
			"string is contained, but if it's exactly the same.")
)

// Print strings properly formatted for the Command Line Interface (CLI) program.
// This is so that when using the CLI program, the user can better understand
// where a string is coming from.
// If you only need to add a newline, don't use this function!
func printCLI(str ...any) {
	fmt.Print("CLI: ")
	fmt.Println(str...)
}

// Same as printCLI(), but only for verbose output and is optional.
func printCLIVerbose(verbose bool, str ...any) {
	if verbose == true {
		printCLI(str...)
	}
}

func printErrInfo(errInfo drugdose.ErrorInfo) {
	if errInfo.Err != nil {
		printCLI(errInfo.Err)
	} else if errInfo.Action != "" {
		printCLI(errInfo.Action)
	}
}

func main() {
	// Initialisation /////////////////////////////////////////////////////
	flag.Parse()

	if flag.NFlag() == 0 {
		printCLI("Try adding -help with space next to the program name! You can read the README file as well.")
		os.Exit(1)
	}

	gotsetcfg := drugdose.InitAllSettings(*sourcecfg, *dbDir, drugdose.DefaultDBName,
		drugdose.DefaultMySQLAccess, *recreateSettings, *recreateSources,
		*verbose, *sourceAddress)

	if *stopOnCfgInit {
		printCLI("Stopping after config file initialization.")
		os.Exit(0)
	}

	ctx, ctx_cancel, err := gotsetcfg.UseConfigTimeout()
	if err != nil {
		printCLI(err)
		os.Exit(1)
	}
	defer ctx_cancel()

	gotsetcfg.InitAllDB(ctx)

	if *stopOnDbInit {
		printCLI("Stopping after database initialization.")
		os.Exit(0)
	}

	db := gotsetcfg.OpenDBConnection(ctx)
	defer db.Close()

	err = gotsetcfg.AddToAllNamesTables(db, ctx, false)
	if err != nil {
		printCLI(err)
		os.Exit(1)
	}
	///////////////////////////////////////////////////////////////////////

	setType := ""
	setValue := ""
	getExact := ""
	if *startTime != "none" {
		setType = drugdose.LogStartTimeCol
		setValue = *startTime
	} else if *endTime != "none" {
		setType = drugdose.LogEndTimeCol
		setValue = *endTime
	} else if *drugname != "none" {
		setType = drugdose.LogDrugNameCol
		setValue = *drugname
	} else if *drugargdose != 0 {
		setType = drugdose.LogDoseCol
		setValue = strconv.FormatFloat(*drugargdose, 'f', -1, 64)
	} else if *drugunits != "none" {
		setType = drugdose.LogDoseUnitsCol
		setValue = *drugunits
	} else if *drugroute != "none" {
		setType = drugdose.LogDrugRouteCol
		setValue = *drugroute
	} else if *drugcost != 0 {
		setType = drugdose.LogCostCol
		setValue = strconv.FormatFloat(*drugcost, 'f', -1, 64)
	} else if *costCur != "" {
		setType = drugdose.LogCostCurrencyCol
		setValue = *costCur
	}

	if *searchExact {
		getExact = setType
		*searchStr = setValue
	}

	if *getPaths {
		printCLI("DB Dir:", gotsetcfg.DBSettings[gotsetcfg.DBDriver].Path)
		err, gotsetdir := drugdose.InitSettingsDir()
		if err != nil {
			printCLI(err)
			os.Exit(1)
		}
		printCLI("Settings Dir:", gotsetdir)
	}

	if *getDBDriver {
		printCLI("Driver:\t", gotsetcfg.DBDriver)
		printCLI("Path:\t", gotsetcfg.DBSettings[gotsetcfg.DBDriver].Path)
		printCLI("Param:\t", gotsetcfg.DBSettings[gotsetcfg.DBDriver].Parameters)
	}

	if *cleanDB {
		err := gotsetcfg.CleanDB(db, ctx)
		if err != nil {
			printCLI(err)
			os.Exit(1)
		}
	}

	if *cleanInfo {
		err := gotsetcfg.CleanInfoTable(db, ctx)
		if err != nil {
			printCLI(err)
			os.Exit(1)
		}
	}

	if *cleanNames {
		err := gotsetcfg.CleanNamesTables(db, ctx, false)
		if err != nil {
			printCLI(err)
			os.Exit(1)
		}
	}

	if *overwriteNames {
		err := gotsetcfg.AddToAllNamesTables(db, ctx, true)
		if err != nil {
			printCLI(err)
			os.Exit(1)
		}
	}

	if *forget {
		gotErrInfo := gotsetcfg.ForgetDosing(db, ctx, nil, *forUser)
		if gotErrInfo.Err != nil {
			printCLI(gotErrInfo.Err)
		} else if gotErrInfo.Err == nil && gotErrInfo.Action != "" {
			printCLI(gotErrInfo.Action)
		}
	}

	remembering := false
	if *drugargdose != 0 && *drugname == "none" && *changeLog == false {
		gotUserLogsErr := gotsetcfg.RecallDosing(db, ctx, nil, *forUser)
		err := gotUserLogsErr.Err
		if err != nil {
			printCLI("Couldn't recall dosing configuration: ", err)
			os.Exit(1)
		} else if gotUserLogsErr.UserLogs != nil {
			remCfg := gotUserLogsErr.UserLogs[0]
			printCLI("Remembering from config using ID:", remCfg.StartTime)
			*forUser = remCfg.Username
			*drugname = remCfg.DrugName
			*drugroute = remCfg.DrugRoute
			*drugunits = remCfg.DoseUnits
			remembering = true
		}
	}

	if *removeInfoDrug != "none" {
		tempErrInfo := gotsetcfg.RemoveSingleDrugInfo(db, ctx, nil, *removeInfoDrug, *forUser)
		printErrInfo(tempErrInfo)
	}

	remAmount := 0
	revRem := false
	if *removeOld != 0 {
		remAmount = *removeOld
	}

	if *removeNew != 0 {
		remAmount = *removeNew
		revRem = true
	}

	if *cleanLogs || remAmount != 0 {
		tempErrInfo := gotsetcfg.RemoveLogs(db, ctx, nil, *forUser,
			remAmount, revRem, *forID, *searchStr, getExact)
		printErrInfo(tempErrInfo)
	}

	inputDose := false
	if *changeLog == false && remembering == false && *getLogs == false &&
		*dontLog == false && *searchExact == false {

		if *drugname != "none" ||
			*drugroute != "none" ||
			*drugargdose != 0 ||
			*drugunits != "none" {

			if *drugname == "none" {
				printCLI("No drug name specified, checkout: gopsydose -help")
			}

			if *drugroute == "none" {
				printCLI("No route specified, checkout: gopsydose -help")
			}

			if *drugargdose == 0 {
				printCLI("No dose specified, checkout: gopsydose -help")
			}

			if *drugunits == "none" {
				printCLI("No units specified, checkout: gopsydose -help")
			}

			if *drugname != "none" && *drugroute != "none" &&
				*drugargdose != 0 && *drugunits != "none" {

				inputDose = true
			}
		}
	}

	if remembering == true {
		inputDose = true
	}

	if inputDose == true || *dontLog == true && *drugname != "none" {
		err, cli := gotsetcfg.InitGraphqlClient()
		fetchErr := false
		if err == nil {
			gotErrInfo := gotsetcfg.FetchFromSource(db, ctx, nil, *drugname, *forUser, cli)
			if gotErrInfo.Err != nil {
				fetchErr = true
				printCLI(gotErrInfo.Err)
			}
		} else {
			printCLI(err)
		}

		if *dontLog == false && fetchErr == false {
			errInfo := gotsetcfg.AddToDoseTable(db, ctx, nil, nil, *forUser, *drugname, *drugroute,
				float32(*drugargdose), *drugunits, float32(*drugperc),
				float32(*drugcost), *costCur, true)
			printErrInfo(errInfo)
		} else if *dontLog == true {
			err, convOutput, convUnit := gotsetcfg.ConvertUnits(db, ctx, *drugname,
				float32(*drugargdose), float32(*drugperc))
			if err != nil {
				printCLI(err)
				os.Exit(1)
			} else {
				convSubs := gotsetcfg.MatchAndReplace(db, ctx, *drugname, drugdose.NameTypeSubstance)
				printCLI(fmt.Sprintf("Didn't log, converted dose: "+
					"%g ; units: %q ; substance: %q ; username: %q",
					convOutput, convUnit, convSubs, *forUser))
			}
		}
	}

	if *changeLog {
		errInfo := gotsetcfg.ChangeUserLog(db, ctx, nil, setType, *forID, *forUser, setValue)
		printErrInfo(errInfo)
	}

	if *dontLog == false && *remember == true {
		errInfo := gotsetcfg.RememberDosing(db, ctx, nil, *forUser, *forID)
		printErrInfo(errInfo)
	}

	// All functions which modify data have finished.

	var gettingLogs bool = false
	var logsLimit bool = false
	var gotUserLogsErr drugdose.UserLogsError

	if *getLogs {
		if *noGetLimit {
			gotUserLogsErr = gotsetcfg.GetLogs(db, ctx, nil, 0, *forID,
				*forUser, false, *searchStr, getExact)
		} else {
			gotUserLogsErr = gotsetcfg.GetLogs(db, ctx, nil, 100, *forID,
				*forUser, false, *searchStr, getExact)
			logsLimit = true
		}
		gettingLogs = true
	} else if *getNewLogs != 0 {
		gotUserLogsErr = gotsetcfg.GetLogs(db, ctx, nil, *getNewLogs, 0,
			*forUser, true, *searchStr, getExact)
		gettingLogs = true
	} else if *getOldLogs != 0 {
		gotUserLogsErr = gotsetcfg.GetLogs(db, ctx, nil, *getOldLogs, 0,
			*forUser, false, *searchStr, getExact)
		gettingLogs = true
	}

	if gettingLogs == true {
		retLogs := gotUserLogsErr.UserLogs
		gotErr := gotUserLogsErr.Err
		if logsLimit == true {
			if retLogs != nil && len(retLogs) == 100 {
				printCLI("By default there is a limit of retrieving " +
					"and printing a maximum of 100 entries. " +
					"To avoid it, use the -no-get-limit option.")
			}
		}

		if gotErr != nil {
			printCLI(gotErr)
		} else {
			gotsetcfg.PrintLogs(retLogs, false)
		}
	}

	if *getLogsCount {
		gotLogCountErr := gotsetcfg.GetLogsCount(db, ctx, *forUser, nil)
		err := gotLogCountErr.Err
		if err != nil {
			printCLI(err)
			os.Exit(1)
		}
		printCLI("Total number of logs:", gotLogCountErr.LogCount, "; for user:", gotLogCountErr.Username)
	}

	if *getTotalCosts {
		gotCostsErr := gotsetcfg.GetTotalCosts(db, ctx, nil, *forUser)
		err := gotCostsErr.Err
		if err != nil {
			printCLI(err)
			os.Exit(1)
		}
		drugdose.PrintTotalCosts(gotCostsErr.Costs, false)
	}

	if *getTimes {
		gotTimeTillErr := gotsetcfg.GetTimes(db, ctx, nil, *forUser, *forID)
		err := gotTimeTillErr.Err
		if err != nil {
			printCLI("Times couldn't be retrieved because of an error:", err)
			os.Exit(1)
		} else {
			err = gotsetcfg.PrintTimeTill(gotTimeTillErr, false)
			if err != nil {
				printCLI("Couldn't print times because of an error:", err)
				os.Exit(1)
			}
		}
	}

	if *getUsers {
		gotAllUsersErr := gotsetcfg.GetUsers(db, ctx, nil, *forUser)
		err = gotAllUsersErr.Err
		ret := gotAllUsersErr.AllUsers
		if err != nil {
			printCLI("Couldn't get users because of an error:", err)
			os.Exit(1)
		} else {
			str := fmt.Sprint("All users: ")
			for i := 0; i < len(ret); i++ {
				str += fmt.Sprintf("%q ; ", ret[i])
			}
			printCLI(str)
		}
	}

	getNamesWhich := "none"
	getNamesValue := ""
	if *getSubNames != "" {
		getNamesWhich = "substance"
		getNamesValue = *getSubNames
	} else if *getRouteNames != "" {
		getNamesWhich = "route"
		getNamesValue = *getRouteNames
	} else if *getUnitsNames != "" {
		getNamesWhich = "units"
		getNamesValue = *getUnitsNames
	}

	if getNamesWhich != "none" {
		gotDrugNamesErr := gotsetcfg.GetAllAltNames(db, ctx, nil,
			getNamesValue, getNamesWhich, false, *forUser)
		subsNames := gotDrugNamesErr.DrugNames
		err = gotDrugNamesErr.Err
		if err != nil {
			printCLI("Couldn't get substance names, because of error:", err)
			os.Exit(1)
		} else {
			fmt.Print("For " + getNamesWhich + ": " + getNamesValue + " ; Alternative names: ")
			for i := 0; i < len(subsNames); i++ {
				fmt.Print(subsNames[i] + ", ")
			}
			fmt.Println()
		}
	}

	getUniqueNames := false
	getInfoNames := false
	useCol := ""
	if *getLocalInfoDrugs {
		getUniqueNames = true
		getInfoNames = true
		useCol = drugdose.InfoDrugNameCol
	} else if *getLoggedDrugs {
		getUniqueNames = true
		useCol = drugdose.LogDrugNameCol
	}

	if getUniqueNames == true {
		gotDrugNamesErr := gotsetcfg.GetLoggedNames(db, ctx, nil, *&getInfoNames, *forUser, useCol)
		err = gotDrugNamesErr.Err
		locinfolist := gotDrugNamesErr.DrugNames
		if err != nil {
			printCLI("Error getting drug names list:", err)
			os.Exit(1)
		} else {
			if getInfoNames {
				str := fmt.Sprint("For source: " + gotsetcfg.UseSource + " ; All local drugs: ")
				for i := 0; i < len(locinfolist); i++ {
					str += fmt.Sprint(locinfolist[i] + " ; ")
				}
				printCLI(str)
			} else {
				str := fmt.Sprint("Logged unique drug names: ")
				for i := 0; i < len(locinfolist); i++ {
					str += fmt.Sprint(locinfolist[i] + " ; ")
				}
				printCLI(str)
			}
		}
	}

	if *getLocalInfoDrug != "none" {
		gotDrugInfoErr := gotsetcfg.GetLocalInfo(db, ctx, nil, *getLocalInfoDrug, *forUser)
		locinfo := gotDrugInfoErr.DrugI
		err = gotDrugInfoErr.Err
		if err != nil {
			printCLI("Couldn't get info for drug because of error:", err)
			os.Exit(1)
		} else {
			gotsetcfg.PrintLocalInfo(locinfo, false)
		}
	}

	if *getDBSize {
		ret := gotsetcfg.GetDBSize()
		retMiB := (ret / 1024) / 1024
		printCLI("Total DB size returned:", retMiB, "MiB ;", ret, "bytes")
	}
}
