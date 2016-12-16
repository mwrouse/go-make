package parser


import (
  "regexp"
  "bufio"
  "os"
  "errors"
  "strings"
  "strconv"
)



// Regular expressions
var varDeclRegex, _ = regexp.Compile("^([A-Za-z0-9]+)(?: )*=(?: )*(.*)?$")
var sectionRegex, _ = regexp.Compile("^([A-Za-z0-9]+):(?: )*$")
var varRegex, _     = regexp.Compile("\\$\\(([A-Za-z0-9]+)\\)")
var commentRegex, _ = regexp.Compile("^( )*#(.*)?$")


/**
 * Name.........: ParseMakefile
 * Parameters...: path (string)- the path to the makefile
 *                wd (string) - current working directory
 * Return.......: map[string]string - map of the global variables in the makefile
 *                map[string][]string - map of the sections to all of the commands in each section
 *                error - any errors
 * Description..: Parses the makefile
 */
func ParseMakefile(path string, wd string) (map[string]string, map[string][]string, error) {
  // Maps for variables and sections
  globalVariables := make(map[string]string)
  localVariables  := make(map[string]string)
  sections        := make(map[string][]string)

  globalVariables["DIR"] = wd // Populate dir global variable

  // Misc variables
  lineCount := ""
  line      := ""
  inSection := false // True when a section has been started, used for detecting global variables
  currentSection := ""

  // Read the file
  contents, err := ReadFile(path)
  if err != nil {
    return nil, nil, err // Pass error to caller
  }


  // Loop through all of the lines in the makefile
  for i, origLine := range contents {
    lineCount = strconv.Itoa(i + 1)
    line = strings.TrimSpace(origLine) // Remove leading and trailing whitespace

    // Parsing
    switch {
    // Comment Line
    case commentRegex.MatchString(line):
      // Do nothing

    // Variable Declaration
    case varDeclRegex.MatchString(line):
        matches := varDeclRegex.FindStringSubmatch(line) // Get the variable name and value
        if len(matches) != 3 {
          return nil, nil, errors.New("Invalid Variable Declaration at line" + lineCount)
        }

        if matches[2] == "" {
          return nil, nil, errors.New("Invalid value for variable " + matches[1] + " on line " + lineCount)
        }

        // Turn variable to all uppercase letters
        matches[1] = strings.ToUpper(matches[1])

        // Add to proper scope
        if inSection {
          // Local variable
          localVariables[matches[1]] = matches[2]
        } else {
          // Global variable
          globalVariables[matches[1]] = matches[2]
        }


    // Section declaration
    case sectionRegex.MatchString(line):
      matches := sectionRegex.FindStringSubmatch(line) // Get matches (section name)

      if len(matches) != 2 {
        return nil, nil, errors.New("Invalid Section Declaration at line " + lineCount)
      }

      // Clear local scope
      localVariables = make(map[string]string)

      // Turn name into uppercase letters
      matches[1] = strings.ToUpper(matches[1])

      if matches[1] == "GLOBAL" {
        inSection = false; // Global section
        currentSection = ""
      } else {
        inSection = true;
        currentSection = matches[1]

        // Add the section
        sections[currentSection] = make([]string, 0)
      }


    // Command (non-blank line that does match the previous cases)
    case line != "":
      // Make sure it is in a section
      if !inSection {
        return nil, nil, errors.New("Command not in a section at line " + lineCount)
      }

      // Line in a section, replace variables with values
      line, err := PopulateVariables(line, lineCount, localVariables, globalVariables)
      if err != nil {
        return nil, nil, err
      }

      // Add line to the current section
      sections[currentSection] = append(sections[currentSection], line)

    } // End switch
  }

  // Return the variables, sections, and no error
  return globalVariables, sections, nil
}



/**
 * Name.........: PopulateVariables
 * Parameters...: line (string) - line to replace variables in
 *                lineCount (string) - the number line in the program
 *                localVariables (map[string]string) - local scope variables
 *                globalVariables (map[string]string) - global scope variables
 * Return.......: string - the new line
 *                error - any errors
 * Description..: Replaces variable reference with variable value
 */
func PopulateVariables(line string, lineCount string, localVariables map[string]string, globalVariables map[string]string) (string, error) {
  // Get matches of all the variables in the string
  matches := varRegex.FindAllString(line, -1)

  // Loop through each variable
  for i, _ := range matches {
    variable := strings.ToUpper(matches[i][2:len(matches[i]) - 1]) // Gets the variable name (removes '$(' and ')' from around the variable name)

    // Check local scope first
    value := localVariables[variable]

    // Check global scope if needed
    if value == ""  {
      value = globalVariables[variable]

      // If still "" or default value return an error
      if value == "" {
        return "", errors.New("Undeclared/Uninitialized variable " + variable + " on line " + lineCount )
      }
    }

    // Replace the variable reference with the value
    line = strings.Replace(line, matches[i], value, -1)
  }

  return line, nil // No errors
}


/**
 * Name.........: ReadFile
 * Parameters...: path (string) - path to the file to read
 * Return.......: []string - array of all the lines in the file
 *                error - any errors
 * Description..: Reads the contents of a file
 */
func ReadFile(path string) ([]string, error) {
  // Check if the file exists
  if !FileExists(path) {
    return nil, errors.New("File " + path + " does not exist")
  }

  f, err := os.Open(path)
  if err != nil {
    return nil, err // Error
  }

  defer f.Close()

  var contents []string
  reader := bufio.NewScanner(f)

  for reader.Scan() {
    contents = append(contents, reader.Text())
  }

  return contents, reader.Err()
}




/**
 * Name.........: FileExists
 * Parameters...: path (string) - path to the file to check
 * Return.......: bool - true if the file exists
 * Description..: Checks if a file exists
 */
func FileExists (name string) bool {
  _, result := os.Stat(name)

  if result != nil {
    if os.IsNotExist(result) {
      return false
    }
  }
  return true
}
