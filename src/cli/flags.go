package cli

import (
	"flag"
	"path/filepath"
)

var usageFlag = flag.Bool("u", false, "Usage")
var allFlag = flag.Bool("a", false, "If this flag is set all collections will be processed")
var saveFlag = flag.Bool("s", false, "If this flag is set the collection will be saved to database")
var interativeFlag = flag.Bool("i", false, "Interative mode")

var wbPath, _ = filepath.Abs("./resources/btf.xlsx")
var filename = flag.String("f", wbPath, "Path to BTF workbook")
var collections = flag.String(
	"c",
	"",
	"Collections separated by commas")

type Flags struct {
	InterativeMode bool
	SaveOperation  bool
	Filename       string
	Collections    string
}

func InitFlags() Flags {
	flag.Parse()

	if *usageFlag {
		panic("Usage:")
	}

	if *allFlag {
		*collections = "user[:10],transaction[:10],service[:10],category"
	}

	if !*allFlag && *collections == "" && !*interativeFlag {
		panic("Inform collections of interest with -c or -a flags.\nUsage:")
	}

	return Flags{
		*interativeFlag,
		*saveFlag,
		*filename,
		*collections}
}

func PrintDefaults() {
	flag.PrintDefaults()
}
