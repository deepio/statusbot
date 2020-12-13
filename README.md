# StatusBot

Very minimalistic script to ping servers and send slack notifications when a status in the server changes. Just build, configure, and run.

### Configure

- To let `statusbot` know what sites to ping, write a YAML file that has the name and the URL of the site you wish to track the status of in Slack. You could replace the `config/test.yaml` file.

```yaml
sites:
  - name: "TestBot"
    url: "http://localhost:8000/"
  - name: "TestBot - Rand"
    url: "http://localhost:8000/rand"
```

- You can write your configuration as a JSON file instead if you want.
  
  > Since [YAML 1.2](https://yaml.org/spec/1.2/spec.html), JSON is regarded as valid YAML. What this means for you the user is that you can choose to write your configuration file in YAML or in JSON and it will not make a difference to `statusbot`.

```JSON
{
  "sites": [
    {
      "name": "TestBot",
      "url": "http://localhost:8000"
    },
    {
      "name": "TestBot - Rand",
      "url": "http://localhost:8000/rand"
    }
  ]
}
```

### Run inside a container

- Set your Slack WebHook in the `docker-compose.yml` file
- `docker-compose up`
  - Or you can also do this:

  ```bash
  docker run \
    -v "/Users/[...]/statusbot/config:/tmp" \
    -e SLACK_WEBHOOK=https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX \
    -t d33pi0/statusbot \
    statusbot # with the channel or a different wait (e.g., -wait 1 -chan "#dev-ops")
  ```

### Build it yourself

- Install [go](https://golang.org/doc/install)
- Run `make build`

### Examples

- From the same directory where the StatusBot binary exists:
`./statusbot -file ./test.yaml`

- You can specify the channel you want to send the messages to:
`./statusbot -chan "#dev-ops"`

- Specify the interval you want the bot to ping the host in seconds.
`./statusbot -wait 1`

- All of the above in random order:
`./statusbot -wait 1 -file ./test.yaml -chan "#dev-ops"`

### TestBot

Use the testbot for testing that the messages are going where you want them to.

- Run the default StatusBot configuration
  - `./statusbot`
- Alternate between bringing up and down the testbot
  - `./testbot`
  - Bring testbot down: <kbd>control</kbd> + <kbd>c</kbd>
  - And back up using command history <kbd>&#8593;</kbd>, then <kbd>Enter</kbd>
