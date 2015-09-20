package serf

import (
  "errors"
  "strings"
  "time"
  "os/exec"

  "github.com/docker/swarm/discovery"
)

// DiscoveryUrl is exported
const DiscoveryURL = "https://discovery.hub.docker.com/v1"

// Discovery is exported
type Discovery struct {
  heartbeat time.Duration
  ttl       time.Duration
  url       string
  token     string
}

func init() {
  Init()
}

// Init is exported
func Init() {
  discovery.Register("serf", &Discovery{})
}

// Initialize is exported
func (s *Discovery) Initialize(ip_port string, heartbeat time.Duration, ttl time.Duration) error {
  s.url = ip_port
  s.heartbeat = heartbeat
  s.ttl = ttl

  // Launch the agent in standby
  err := exec.Command("./launch_agent.sh").Run()
  if err != nil {
    return err
  }

  return nil
}

// Fetch returns the list of entries for the discovery service at the specified endpoint
func (s *Discovery) fetch() (discovery.Entries, error) {
  // Here, we contact the serf agent and ask for a list of members
  // We have to send back entries, which are basically a host and a port (see discovery.go)
  output, err := exec.Command("./agent_members.sh").Output()

  if err != nil {
    return nil, err
  }


  lines := strings.Split(string(output[:]), "\n")
  var addrs []string
  for _, line := range lines {
    fields := strings.Fields(line)
    if len(fields) != 3 {
      return nil, errors.New("Error while parsing the output of serf members : Wrong number of fields")
    }
    addrs = append(addrs, fields[1])
  }

  return discovery.CreateEntries(addrs)
}

// Watch is exported
func (s *Discovery) Watch(stopCh <-chan struct{}) (<-chan discovery.Entries, <-chan error) {
  ch := make(chan discovery.Entries)
  ticker := time.NewTicker(s.heartbeat)
  errCh := make(chan error)

  go func() {
    defer close(ch)
    defer close(errCh)

    // Send the initial entries if available.
    currentEntries, err := s.fetch()
    if err != nil {
      errCh <- err
    } else {
      ch <- currentEntries
    }

    // Periodically send updates.
    for {
      select {
      case <-ticker.C:
        newEntries, err := s.fetch()
        if err != nil {
          errCh <- err
          continue
        }

        // Check if the file has really changed.
        if !newEntries.Equals(currentEntries) {
          ch <- newEntries
        }
        currentEntries = newEntries
      case <-stopCh:
        ticker.Stop()
        return
      }
    }
  }()

  return ch, errCh
}

// Register adds a new entry identified by the into the discovery service
func (s *Discovery) Register(addr string) error {
  // Launch the agent in standby
  err := exec.Command("./agent_join.sh", s.url).Run()
  if err != nil {
    return err
  }

  return nil
}
