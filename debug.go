package simplessh

import (
	"fmt"
	"log"
)

//--- Send messages to the logging subsystem, but only when the provided level
// is equal-to or lesser-than the global DebugLevel.
func Logit(msgLevel int, msg string) {
  if DebugLevel >= msgLevel {
    log.Printf("[%d] %s", msgLevel, msg)
  }
}

//--- Special func to Log Config info defined in a Client.
func (self *Client) LogClientConfig() {

	//- Logging data in client.Opts requires DebugLevel > 3.
	if DebugLevel >= 4 {
		Logit (4, "-------------------------------------")
		Logit (4, "client.Opts Info:")
		Logit (4, "client.Opts.Address = "+self.Opts.Address)
		Logit (4, "client.Opts.AuthMethod = "+self.Opts.AuthMethod)
		Logit (4, fmt.Sprintf("client.Opts.DebugLevel = %d", self.Opts.DebugLevel))
		Logit (4, fmt.Sprintf("client.Opts.IgnoreHostKey = %d", self.Opts.IgnoreHostKey))
		Logit (4, "client.Opts.SuccessPattern = "+self.Opts.SuccessPattern)
		Logit (4, fmt.Sprintf("client.Opts.TerminalColumns = %d", self.Opts.TerminalColumns))
		Logit (4, fmt.Sprintf("client.Opts.TerminalEcho = %d", self.Opts.TerminalEcho))
		Logit (4, "client.Opts.TerminalMode = "+self.Opts.TerminalMode)
		Logit (4, fmt.Sprintf("client.Opts.TerminalRows = %d", self.Opts.TerminalRows))
		Logit (4, fmt.Sprintf("client.Opts.Timeout = %d", self.Opts.Timeout))
		Logit (4, "client.Opts.User = "+self.Opts.User)
		Logit (4, fmt.Sprintf("client.Opts.UsePty = %d", self.Opts.UsePty))
		Logit (4, "-------------------------------------")
	}
}

//--- Special func to Log Config info defined in a Command.
func (self *Command) LogCommandConfig() {
	//- Logging data in Command requires DebugLevel > 2.
	if DebugLevel >= 3 {
		Logit (3, "-------------------------------------")
		Logit (3, "Command Info:")
		Logit (3, "Command.Exec = "+self.Exec)
		for desc, rgx := range self.ValidPatterns {
			Logit (3, fmt.Sprintf("Command.ValidPatterns.%s = %q", desc, rgx))
		}
		for desc, rgx := range self.ErrorPatterns {
			Logit (3, fmt.Sprintf("Command.ErrorPatterns.%s = %q", desc, rgx))
		}
	}

	//- Logging the compiled Regexp data requires DebugLevel > 4.
	if DebugLevel >= 5 {
		Logit (5, "Command compiled Regex values (Should match the Command.XxxPatterns values):")
		for desc, rgx := range self.ValidRegex {
			Logit (5, fmt.Sprintf("  Command.ValidRegex.%s = %q", desc, rgx.String()))
		}
		for desc, rgx := range self.ErrorRegex {
			Logit (5, fmt.Sprintf("  Command.ErrorRegex.%s = %q", desc, rgx.String()))
		}
	}
	Logit (3, "-------------------------------------")
}
