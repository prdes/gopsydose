package drugdose

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/hasura/go-graphql-client"

	"database/sql"

	// MySQL driver needed for sql module
	_ "github.com/go-sql-driver/mysql"
	// SQLite driver needed for sql module
	_ "github.com/glebarez/go-sqlite"
)

type PsychonautwikiSubstance []struct {
	Name string

	Roas []struct {
		Name string

		Dose struct {
			Units     string
			Threshold float64
			Light     struct {
				Min float64
				Max float64
			}
			Common struct {
				Min float64
				Max float64
			}
			Strong struct {
				Min float64
				Max float64
			}
		}

		Duration struct {
			Onset struct {
				Min   float64
				Max   float64
				Units string
			}

			Comeup struct {
				Min   float64
				Max   float64
				Units string
			}

			Peak struct {
				Min   float64
				Max   float64
				Units string
			}

			Offset struct {
				Min   float64
				Max   float64
				Units string
			}

			Total struct {
				Min   float64
				Max   float64
				Units string
			}
		}
	}
}

// Used to initialise the GraphQL struct, using the source address from
// the drugdose Config struct.
//
// returns the GraphQL struct used with github.com/hasura/go-graphql-client
func (cfg Config) InitGraphqlClient() (error, graphql.Client) {
	const printN string = "InitGraphqlClient()"

	client := graphql.Client{}

	if !cfg.AutoFetch {
		printNameVerbose(cfg.VerbosePrinting, printN, "Automatic fetching is disabled, returning.")
		return nil, client
	}

	var proxy func(*http.Request) (*url.URL, error) = nil
	if cfg.ProxyURL != "" && cfg.ProxyURL != "none" {
		goturl, err := url.Parse(cfg.ProxyURL)
		if err != nil {
			err = errors.New(sprintName(printN, err))
			return err, client
		}
		proxy = http.ProxyURL(goturl)
	}
	var CustomTransport http.RoundTripper = &http.Transport{
		Proxy:                 proxy,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	httpClient := http.Client{
		Transport: CustomTransport,
	}
	gotsrcData := GetSourceData()
	if gotsrcData == nil {
		return errors.New(sprintName(printN, "GetSourceData() returned nil, returning.")), client
	}
	api := gotsrcData[cfg.UseSource].API_ADDRESS
	apiURL := "https://" + api
	client_new := graphql.NewClient(apiURL, &httpClient)
	return nil, *client_new
}

// FetchPsyWiki gets information from Psychonautwiki about a given substance
// and stores it in the local info table. The table is determined by the
// source chosen in the Config struct. The name of the table is the same as the
// name of the source, in this case "psychonautwiki".
//
// This function is meant to be run concurrently.
//
// db - open database connection
//
// ctx - context to be passed to sql queries
//
// errChannel - the gorouting channel which returns the errors
//
// drugname - the substance to get information about
//
// client - the initialised structure for the graphql client,
// best done using InitGraphqlClient(), but can be done manually if needed
func (cfg Config) FetchPsyWiki(db *sql.DB, ctx context.Context,
	errChannel chan error, drugname string, client graphql.Client) {
	const printN string = "FetchPsyWiki()"

	if !cfg.AutoFetch {
		printNameVerbose(cfg.VerbosePrinting, printN, "Automatic fetching is disabled, returning.")
		errChannel <- nil
		return
	}

	drugname = cfg.MatchAndReplace(db, ctx, drugname, "substance")

	ret := checkIfExistsDB(db, ctx,
		"drugName",
		"psychonautwiki",
		cfg.DBDriver,
		cfg.DBSettings[cfg.DBDriver].Path,
		nil,
		drugname)
	if ret {
		printNameVerbose(cfg.VerbosePrinting, printN, "Drug already in DB, returning. No need to fetch anything from Psychonautwiki.")
		errChannel <- nil
		return
	}

	printName(printN, "Fetching from source:", cfg.UseSource)

	// This is the graphql query for Psychonautwiki.
	// The way it works is, the full query is generated
	// using the PsychonautwikiSubstance struct.
	var query struct {
		PsychonautwikiSubstance `graphql:"substances(query: $dn)"`
	}

	// Since the query has to be a string, the module has provided
	// an argument allowing to map a variable to the string.
	variables := map[string]interface{}{
		"dn": drugname,
	}

	err := client.Query(ctx, &query, variables)
	if err != nil {
		errChannel <- errors.New(sprintName(printN, "Error from Psychonautwiki API: ", err))
		return
	}

	InfoDrug := []DrugInfo{}

	if len(query.PsychonautwikiSubstance) != 0 {
		subs := query.PsychonautwikiSubstance
		for i := 0; i < len(subs); i++ {
			if len(subs[i].Roas) != 0 {
				for o := 0; o < len(subs[i].Roas); o++ {
					printNameVerbose(cfg.VerbosePrinting, printN, "From source:", cfg.UseSource, "; Substance:", subs[i].Name,
						"; Route:", subs[i].Roas[o])

					tempInfoDrug := DrugInfo{}

					tempInfoDrug.DrugName = subs[i].Name
					tempInfoDrug.DrugRoute = subs[i].Roas[o].Name
					tempInfoDrug.Threshold = float32(subs[i].Roas[o].Dose.Threshold)
					tempInfoDrug.LowDoseMin = float32(subs[i].Roas[o].Dose.Light.Min)
					tempInfoDrug.LowDoseMax = float32(subs[i].Roas[o].Dose.Light.Max)
					tempInfoDrug.MediumDoseMin = float32(subs[i].Roas[o].Dose.Common.Min)
					tempInfoDrug.MediumDoseMax = float32(subs[i].Roas[o].Dose.Common.Max)
					tempInfoDrug.HighDoseMin = float32(subs[i].Roas[o].Dose.Strong.Min)
					tempInfoDrug.HighDoseMax = float32(subs[i].Roas[o].Dose.Strong.Max)
					tempInfoDrug.DoseUnits = subs[i].Roas[o].Dose.Units
					tempInfoDrug.OnsetMin = float32(subs[i].Roas[o].Duration.Onset.Min)
					tempInfoDrug.OnsetMax = float32(subs[i].Roas[o].Duration.Onset.Max)
					tempInfoDrug.OnsetUnits = subs[i].Roas[o].Duration.Onset.Units
					tempInfoDrug.ComeUpMin = float32(subs[i].Roas[o].Duration.Comeup.Min)
					tempInfoDrug.ComeUpMax = float32(subs[i].Roas[o].Duration.Comeup.Max)
					tempInfoDrug.ComeUpUnits = subs[i].Roas[o].Duration.Comeup.Units
					tempInfoDrug.PeakMin = float32(subs[i].Roas[o].Duration.Peak.Min)
					tempInfoDrug.PeakMax = float32(subs[i].Roas[o].Duration.Peak.Max)
					tempInfoDrug.PeakUnits = subs[i].Roas[o].Duration.Peak.Units
					tempInfoDrug.OffsetMin = float32(subs[i].Roas[o].Duration.Offset.Min)
					tempInfoDrug.OffsetMax = float32(subs[i].Roas[o].Duration.Offset.Max)
					tempInfoDrug.OffsetUnits = subs[i].Roas[o].Duration.Offset.Units
					tempInfoDrug.TotalDurMin = float32(subs[i].Roas[o].Duration.Total.Min)
					tempInfoDrug.TotalDurMax = float32(subs[i].Roas[o].Duration.Total.Max)
					tempInfoDrug.TotalDurUnits = subs[i].Roas[o].Duration.Total.Units

					InfoDrug = append(InfoDrug, tempInfoDrug)
				}
			} else {
				printName(printN, "No roas for:", subs[i])
			}
		}

		if len(InfoDrug) != 0 {
			errChannel2 := make(chan error)
			go cfg.AddToInfoDB(db, ctx, errChannel2, InfoDrug)
			err := <-errChannel2
			if err != nil {
				errChannel <- errors.New(sprintName(printN, "Data couldn't be added to info DB, because of an error: ", err))
				return
			}
			printName(printN, "Data added to info DB successfully.")
		} else {
			errChannel <- errors.New(sprintName(printN, "Struct array is empty, nothing added to DB."))
			return
		}
	} else {
		errChannel <- errors.New(sprintName(printN, "The Psychonautwiki API returned nothing."))
		return
	}

	errChannel <- nil
}
