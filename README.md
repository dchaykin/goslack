# goslack
Slack notifier for go (Package)

### Using
import "dchaykin/goslack"

```
func main() {
  ci := goslack.ConfigItem{ Level: "ERROR", URL: "https://hooks.slack.com/services/TT12345/B012345/ABCD123456" }
  if err := goslack.AddConfig(ci); err != nil {
    fmt.Fatal(err)
  }
  goslack.Errorf("This is an error sent to slack")
}
```
