package simplessh

import (
  "bufio"
  "golang.org/x/crypto/ssh"
  "io"
  "regexp"
	"time"
)

//--- These are defaults to some Opts, so we don't need to define them
// every time.
const (
  DefaultTerminalColumns  = 80
  DefaultTerminalRows     = 40
  DefaultTerminalEcho     = 1
  DefaultTerminalMode     = `vt100`
	//- Should cover most prompts.
  DefaultSuccessPattern   = `[\#\$\>]\s*$`
	//- Attempt to connnect for 60 seconds before giving up.
	DefaultTimeout					=	time.Second * 60
	DefaultDebugLevel				= 1
)

//--- SSH Client
type Client struct {
  Session     *ssh.Session
  Config      *ssh.ClientConfig
  Host        string
  Port        int
  StdinPipe   io.WriteCloser
  StdoutPipe  *bufio.Reader
  Opts        *Opts
}

//--- Options for SSH Client
type Opts struct {
  Port              int
  User              string
  Password          string
  AuthMethod        string
  IgnoreHostKey     bool
  Address           string
  UsePty            bool
  TerminalEcho      uint32
  TerminalColumns   int
  TerminalRows      int
  TerminalMode      string
  //- REGEX used for determining that we established a successful Session.
  SuccessPattern    string
	Timeout						time.Duration
	DebugLevel				int
}

//--- SSH Command
type Command struct {
  Exec          string
  ValidPatterns map[string]string
  ErrorPatterns map[string]string

	//- Compiled forms of the above Patterns. These are defined when RunCommand
	// is called for the first time.
	//
	// Compiled form of ValidPatterns.
  ValidRegex    map[string]*regexp.Regexp
	// Compiled form of ErrorPatterns
  ErrorRegex    map[string]*regexp.Regexp
}

//--- Make the Debug Level globally available.
var DebugLevel int
