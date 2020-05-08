# goslack
Slack notifier for go (Package)

### Using

```
import "dchaykin/goslack"

func main() {
  ciError := goslack.ConfigItem{ Level: "ERROR", URL: "https://hooks.slack.com/services/TT12345/B012345/ERROR123456" }
  ciWarning := goslack.ConfigItem{ Level: "WARNING", URL: "https://hooks.slack.com/services/TT12345/B012345/WARNING123456" }
  ciInfo := goslack.ConfigItem{ Level: "INFO", URL: "https://hooks.slack.com/services/TT12345/B012345/INFO123456" }
  
  if err := goslack.AddConfig(ciError, ciWarning, ciInfo); err != nil {
    fmt.Fatal(err)
  }
  
  goslack.Errorf("This is an error")
  goslack.Warningf("This is a warning")
  goslack.Infof("This is a piece of information")
}
```
