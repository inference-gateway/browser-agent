# Example Playwright Automation Script

This script demonstrates how to use the Playwright automation framework to perform basic browser actions such as navigating to a webpage, filling out a form, and taking a screenshot.


Configure the environment variables as needed:

```bash
cp .env.example .env
```

** Add at least two providers, in this example Google and DeepSeek.

First bring up all the containers:

```bash
docker compose up --build
```

Go into the CLI for convenience:

```bash
docker compose run --rm cli
```

Ask the following:

```text
Please visit http://demo-site which is running locally and take a screenshot of the homepage. Use the agent.
```

You would see the CLI (A2A agent client) submitting a task to the A2A agent server and the screenshot will appear in the `screenshots` directory since it's mounted as a volume.

```text
Please visit http://demo-site which is running locally and collect all of the prices, write them to a CSV file. Use the agent.
```

You would see the CLI (A2A agent client) submitting a task to the A2A agent server and the csv file with all of the prices of the website will appear inside of the artifacts directory.

Check the logs to see that the browser indeed went to the demo site and took a screenshot:

```bash
docker compose logs -f demo-site
```

Also you can check the task was successfully submitted to the agent and it's available using the a2a debugger:

```bash
docker compose run --rm a2a-debugger tasks list
```

Finally clean up:

```bash
docker compose down
```
