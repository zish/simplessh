package simplessh

import (
  "bufio"
  "fmt"
  "golang.org/x/crypto/ssh"
  "golang.org/x/crypto/ssh/agent"
  "io/ioutil"
  "net"
  "os"
  "regexp"
)

//--- Generate a new SSH Client.
func New(opts *Opts) (*Client, error) {
  //- Set some Option Defaults.
  if opts.UsePty == true {
    if opts.TerminalColumns == 0 {
      opts.TerminalColumns = DefaultTerminalColumns
    }
    if opts.TerminalRows == 0 {
      opts.TerminalRows = DefaultTerminalRows
    }
    if opts.TerminalEcho == 0 {
      opts.TerminalEcho = DefaultTerminalEcho
    }
    if opts.TerminalMode == "" {
      opts.TerminalMode = DefaultTerminalMode
    }
  }

  if opts.SuccessPattern == "" {
    opts.SuccessPattern = DefaultSuccessPattern
  }

  if opts.Timeout < 1 {
    opts.Timeout = DefaultTimeout
  }

  if opts.DebugLevel < 1 {
    opts.DebugLevel = DefaultDebugLevel
  }

	//- Globally define the Debug Level.
	DebugLevel = opts.DebugLevel
	Logit(3, fmt.Sprintf("Debug Logging Level set to %d", DebugLevel))

  clientConfig := &ssh.ClientConfig{
    User: opts.User,
		Timeout: opts.Timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
  }

  //- Define SSH Host Key callback mechanism.
  if opts.IgnoreHostKey == true {
    clientConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

  //- Define how SSH authentication will be handled.
  if opts.AuthMethod == "agent" {
    clientConfig.Auth = []ssh.AuthMethod{
      SSHAgent(),
    }
  } else if opts.AuthMethod == "password" {
    clientConfig.Auth = []ssh.AuthMethod {
      ssh.Password(opts.Password),
    }
  } else {
    return nil, fmt.Errorf("Invalid SSH Authentication Method: %q", opts.AuthMethod)
  }

  client := &Client{
    Config:  clientConfig,
    Host:    opts.Address,
    Port:    opts.Port,
    Opts:    opts,
  }

	client.LogClientConfig() //- DEBUG

  return client, nil
}

//--- Establish the SSH connection, and get a usable Session.
func (self *Client) startSession() error {
  var (
    err        error
    connection *ssh.Client
  )

  //- Establish the SSH connection.
  if connection, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", self.Host, self.Port), self.Config); err != nil {
    return err
  }

  //- Initialize the connection; Get it ready for sending/receiving data.
  if self.Session, err = connection.NewSession(); err != nil {
    return err
  }

  //- Determine how we should set up the terminal (Use a PTY, echo, etc.).
  if self.Opts.UsePty == true {
    //- enable/disable terminal echo.
    modes := ssh.TerminalModes{
      ssh.ECHO: self.Opts.TerminalEcho,
    }

    if err := self.Session.RequestPty(self.Opts.TerminalMode,
                                      self.Opts.TerminalColumns,
                                      self.Opts.TerminalRows,
                                      modes); err != nil {

      //- Return both errors, if we have problems with closing the Session.
      if closeErr := self.Session.Close(); closeErr != nil {
        errStr := fmt.Sprintf("%s\n%s", err, closeErr)
        return fmt.Errorf(errStr)
      }

      return err
    }
  }

  outPipe, _ := self.Session.StdoutPipe()
  self.StdoutPipe = bufio.NewReader(outPipe)

  if self.StdinPipe, err = self.Session.StdinPipe(); err != nil {
    return err

  } else {

    if err = self.Session.Shell(); err != nil {
      return err
    }
  }

  //- Loop until we get a match against self.Opts.SuccessPattern.
  // Determine if we received what we expected to see after login.
  pattern := regexp.MustCompile(self.Opts.SuccessPattern)
  var output string
  for 1 == 1 {
    b, _ := self.StdoutPipe.ReadByte()
    output = output + string(b)

    if r := pattern.FindStringSubmatch(output); r != nil {
      return nil
    }
  }

  return nil
}

//--- SSH Public key authentication, for a specific Public key (WIP).
//noinspection GoUnusedExportedFunction
func PublicKeyFile(file string) ssh.AuthMethod {
  buffer, err := ioutil.ReadFile(file)
  if err != nil {
    return nil
  }

  key, err := ssh.ParsePrivateKey(buffer)
  if err != nil {
    return nil
  }

  return ssh.PublicKeys(key)
}

//--- Use an available SSH Agent.
// ** Env var SSH_AUTH_SOCK must be set.
func SSHAgent() ssh.AuthMethod {
  if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
    return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
  }

  return nil
}

func (self *Client) Connect() error {
  return self.startSession()
}

//--- Disconnect and close the SSH Session.
func (self *Client) Disconnect() error {
  if self.Session == nil {
    return fmt.Errorf("SSH Session is not initialized.")
  }

  //- Make sure we successfully disconnect.
  if err := self.Session.Close(); err != nil {
    return err
  }

  return nil
}
