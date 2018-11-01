package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/integrii/flaggy"
)

func parseArgs() {
	// Set your program's name and description.  These appear in help output.
	// flaggy.SetName(color.CyanString("%s %s", appname, version))
	flaggy.SetName(appname)
	flaggy.SetDescription("Eat my shorts to be best at sports (" + version + ") ")
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://bitbucket.org/briiC/findthese/\n"
	flaggy.DefaultParser.AdditionalHelpPrepend += strings.Repeat(".", 80)

	// add a global bool flag for fun
	flaggy.String(&argSourcePath, "s", "src", "Source path of directory -- REQUIRED")
	flaggy.String(&argEndpoint, "u", "url", "URL endpoint to hit -- REQUIRED")
	flaggy.String(&argMethod, "m", "method", "HTTP Method to use (default: "+argMethod+")")
	flaggy.String(&argOutput, "o", "output", "Output report to file (default: "+argOutput+")")
	flaggy.String(&argOutput, "z", "delay", "Delay every request for N milliseconds (default: "+fmt.Sprintf("%d", argDelay)+")")
	flaggy.StringSlice(&argSkip, "", "skip", "Skip files with these extensions (default: "+fmt.Sprintf("%v", argSkip)+")")
	flaggy.StringSlice(&argSkipExts, "", "skip-ext", "Skip files with these extensions (default: "+fmt.Sprintf("%v", argSkipExts)+")")
	flaggy.Bool(&argDirOnly, "", "dir-only", "Scan directories only")

	// set the version and parse all inputs into variables
	flaggy.SetVersion(version)
	flaggy.Parse()

	// On missing params show help
	if argSourcePath == "" || argEndpoint == "" {
		flaggy.ShowHelpAndExit("")
	}

	// Validate
	if err := validateArgs(); err != nil {
		color.Red("\n%v\n\n", err)
		return
	}

}

// Validate arguments
func validateArgs() error {

	// Does source path exists
	if _, err := os.Stat(argSourcePath); os.IsNotExist(err) {
		// path/to/whatever does not exist
		return fmt.Errorf("Source path [-s, --src]: \n\t%v", err)
	}

	// NB! Do not check here if URL is available!
	// Because of different configurations given base URL could not be "200 OK"
	// Also there could be configurations where only valid files gives different response and others fails

	// // Trailing slash - URL must end with slash
	// argEndpoint = strings.TrimSuffix(argEndpoint, "/") + "/"

	// Method uppercase - necessary only for visual appearance
	argMethod = strings.ToUpper(argMethod)

	// Delay
	argDelay = int(math.Abs(float64(argDelay)))

	// Skpi files/dirs
	argSkip = normalizeArgSlice(argSkip)

	// Skiped extensions
	var exts []string
	argSkipExts = normalizeArgSlice(argSkipExts)
	for _, ext := range argSkipExts {
		ext = strings.Trim(ext, " .")
		ext = strings.ToLower(ext)
		ext = "." + ext // must be prefixed with dot
		if ext != "" {
			exts = append(exts, ext)
		}
	}
	argSkipExts = exts

	// No errors
	return nil
}

func printUsedArgs() {
	fmt.Println(strings.Repeat("-", 80))
	color.Cyan("%20s: %s", "Source path", argSourcePath)
	color.Cyan("%20s: %s", "URL", argEndpoint)
	color.Cyan("%20s: %s", "Method", argMethod)
	color.Cyan("%20s: %v", "Dir only", argDirOnly)
	color.Cyan("%20s: %s", "Output", argOutput)
	color.Cyan("%20s: %d (ms)", "Delay", argDelay)
	color.Cyan("%20s: %v", "Ignore dir/files", argSkip)
	color.Cyan("%20s: %v", "Ignore extensions", argSkipExts)
	color.Cyan("%20s: %v", "Mutation options", argBackups)
	fmt.Println(strings.Repeat("-", 80))
}

func normalizeArgSlice(arr []string) []string {
	s := strings.Join(arr, ",")

	// all to one separator
	s = strings.Replace(s, ";", ",", -1)
	s = strings.Replace(s, "/", ",", -1)
	s = strings.Replace(s, "|", ",", -1)

	// back to slice and items that added in as cli also now separated
	arr = strings.Split(s, ",")
	return arr
}