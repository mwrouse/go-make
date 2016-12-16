package main

import (
  "os"
  "os/exec"
  "strings"
  "github.com/fatih/color"
  "errors"
  "go-make/parser"
  "flag"
)


/**
 * Name.........: main
 * Description..: Runs a makefile
 */
func main() {
  wd, err := os.Getwd()
  HandleError(err)

  // Get the flags
  makefileFlag := flag.String("f", "makefile", "The makefile to execute")
  sectionFlag := flag.String("s", "ALL", "Section of makefile to run")
  flag.Parse()

  makefile := wd + "\\" + *makefileFlag
  sectionToRun := strings.ToUpper(*sectionFlag)

  _, sections, err := parser.ParseMakefile(makefile, wd)
  if err != nil {
    HandleError(err)
  }

  // Make sure section is populated
  if len(sections[sectionToRun]) == 0 {
    ShowError("Invalid section name or no commands in the section")
  }

  // Run all the commands in the start section
  ExecSection(sectionToRun, sections)

  // Finish
  color.Green("Make finished")
}




/**
 * Name.........: ExecSection
 * Parameters...: section (string) - the name of the section to execute
 *                sections (map[string][]string) - map of all of the sections
 * Description..: Runs all the commands in a section
 */
func ExecSection(section string, sections map[string][]string) {
  headline := color.New(color.FgYellow, color.Bold)

  // Loop through all the commands in the section
  for i, _ := range sections[section] {
    command := sections[section][i]

    // TODO: Check if the command is calling another section
    if len(sections[strings.ToUpper(command)]) != 0 {
      // Execute the commands in the section
      ExecSection(strings.ToUpper(command), sections)

    } else {
      // Execute the command
      cmd := exec.Command("cmd", "/C", command)
      output, err := cmd.CombinedOutput()

      // Show results of the command
      headline.Println(command)
      if err != nil {
        color.Red("\t%s", string(output)) // Command had an error
      } else if string(output) != ""{
        color.White("\t%s", string(output)) // No error (maybe)
      }
    }
  }

} // End ExecSection



// Handles an error if one occurs
func HandleError(err error) {
  if err != nil {
    color.Red("Error: %s", err)
    os.Exit(1)
  }
}

func ShowError(params ...string) {
  err := ""
  for i := range params {
    err = err + params[i]
  }
  HandleError(errors.New(err))
}
