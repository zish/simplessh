package simplessh

import (
  "fmt"
  "regexp"
)

//--- Compile Valid and Error Command Regex patterns
// (so we don't need to every time they are used).
//
// They are added to cmd as cmd.ValidRegex and cmd.ErrorRegex.
func (cmd *Command) CompileCmdRegex() {
  if len(cmd.ValidRegex) == 0 {
    cmd.ValidRegex = make(map[string]*regexp.Regexp)

    for desc, pattern := range cmd.ValidPatterns {
      cmd.ValidRegex[desc] = regexp.MustCompile(pattern)
    }
  }

  if len(cmd.ErrorRegex) == 0 {
    cmd.ErrorRegex = make(map[string]*regexp.Regexp)

    for desc, pattern := range cmd.ErrorPatterns {
      cmd.ErrorRegex[desc] = regexp.MustCompile(pattern)
    }
  }
}

//--- Send input to a remote SSH Session. Return a string containing the output.
func (self *Client) RunCommand(cmd *Command) (string, []error) {
  var (
    output string
    errs   []error
  )

  //- Compile our "Valid" and "Error" REGEX patterns, if not done already.
  // This will allow us to reuse it with compiling it every time.
  cmd.CompileCmdRegex()

	cmd.LogCommandConfig() //- DEBUG

  //- "Write" the command to the remote SSH session.
  self.StdinPipe.Write([]byte(cmd.Exec))

  //- Loop until we've gotten all output from the remote SSH Session.
  done := false
  lastError := ""
  for done != true {
    b, _ := self.StdoutPipe.ReadByte()
    output = output + string(b)

    //- Check for errors, if we have any ErrorPatterns to test.
    if len(cmd.ErrorPatterns) > 0 {
      if matchedError, errDesc := testRegex(output, cmd.ErrorRegex); matchedError == true {
        if lastError != errDesc {
          errs = append(errs, fmt.Errorf("Matched error pattern: %s", errDesc))
          lastError = errDesc
        }
      }
    }

    //- Check for Valid output. Continue retrieving bytes until we see a
    // "valid" pattern.
    if done, _ = testRegex(output, cmd.ValidRegex); done == true {
      //- Make sure there isn't any more left to read.
//			time.Sleep(time.Second)
      if buff := self.StdoutPipe.Buffered(); buff != 0 {
        done = false
      }
    }
  }

  return output, errs
}

//--- Test Regexp maps.
func testRegex(input string, regexMap map[string]*regexp.Regexp) (bool, string) {
  var outputDescription string
  for desc, regex := range regexMap {

    outputDescription = desc
    if r := regex.FindStringSubmatch(input); r != nil {
      return true, outputDescription
    }
  }

  return false, outputDescription
}
