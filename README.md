# StatusBot

Very minimalistic script to ping servers and send slack notifications when a status in the server changes. Just build, configure, and run.

### Build

- Install go
- `make build`

### Configure

- Write a JSON file that has the name and the URL of the site you wish to track the status of in Slack.

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

- Set your Slack WebHook as an environment variable:
  - `export SLACK_WEBHOOK=https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX`

### Execute

- From the same directory where the StatusBot binary exists:
`./statusbot -file ./test.json`

- You can specify the channel you want to send the messages to:
`./statusbot -chan "#dev-ops" -file ./test.json`

- Specify the interval you want the bot to ping the host in seconds.
`./statusbot -file ./test.json -wait 1`

- All of the above in random order:
`./statusbot -wait 1 -file ./test.json -chan "#dev-ops"`

### TestBot

Use the testbot for testing that the messages are going where you want them to.

- Run the default StatusBot configuration
  - `./statusbot`
- Alternate between bringing up and down the testbot
  - `./testbot`
  - Bring testbot down: <kbd>control</kbd> + <kbd>c</kbd>
  - And back up using command history <kbd>&#8593;</kbd>, then <kbd>Enter</kbd>
