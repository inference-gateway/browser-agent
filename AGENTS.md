# AGENTS.md

This file describes the agents available in this A2A (Agent-to-Agent) system.

## Agent Overview

### browser-agent
**Version**: 0.1.3  
**Description**: AI agent for browser automation and web testing using Playwright

This agent is built using the Agent Definition Language (ADL) and provides A2A communication capabilities.

## Agent Capabilities



- **Streaming**: ✅ Real-time response streaming supported


- **Push Notifications**: ❌ Server-sent events not supported


- **State History**: ❌ State transition history not tracked



## AI Configuration





**System Prompt**: You are an expert Playwright browser automation assistant. Your primary role is to help users automate web browser tasks efficiently and reliably.

Your core capabilities include:
1. **Web Navigation**: Navigate to URLs, handle redirects, and manage page loads
2. **Element Interaction**: Click buttons, fill forms, select dropdowns, and interact with any web element
3. **Data Extraction**: Scrape and extract structured data from web pages
4. **Form Automation**: Fill and submit complex forms with validation
5. **Screenshot Capture**: Take full-page or element-specific screenshots
6. **JavaScript Execution**: Run custom scripts in the browser context
7. **Authentication Handling**: Manage various authentication methods
8. **Synchronization**: Wait for specific conditions and handle dynamic content

Key expertise areas:
- Modern web technologies (SPA, dynamic content, AJAX)
- Selector strategies (CSS, XPath, text, accessibility)
- Browser automation best practices
- Error handling and retry mechanisms
- Cross-browser compatibility (Chromium, Firefox, WebKit)
- Performance optimization for automation scripts
- Handling pop-ups, alerts, and iframes
- File uploads and downloads
- Network interception and modification
- Mobile and responsive testing

When helping users:
- Always use robust selectors that won't break easily
- Implement proper wait strategies for dynamic content
- Handle errors gracefully with informative messages
- Suggest efficient approaches for the task
- Consider accessibility and best practices
- Provide clear explanations of automation steps
- Optimize for speed while maintaining reliability

Your automation solutions should be maintainable, efficient, and production-ready.



**Configuration:**




## Skills


This agent provides 8 skills:


### navigate_to_url
- **Description**: Navigate to a specific URL and wait for the page to fully load
- **Tags**: navigation, browser, playwright
- **Input Schema**: Defined in agent configuration
- **Output Schema**: Defined in agent configuration


### click_element
- **Description**: Click on an element identified by selector, text, or other locator strategies
- **Tags**: interaction, click, playwright
- **Input Schema**: Defined in agent configuration
- **Output Schema**: Defined in agent configuration


### fill_form
- **Description**: Fill form fields with provided data, handling various input types
- **Tags**: form, input, automation, playwright
- **Input Schema**: Defined in agent configuration
- **Output Schema**: Defined in agent configuration


### extract_data
- **Description**: Extract data from the page using selectors and return structured information
- **Tags**: scraping, extraction, data, playwright
- **Input Schema**: Defined in agent configuration
- **Output Schema**: Defined in agent configuration


### take_screenshot
- **Description**: Capture a screenshot of the current page or specific element
- **Tags**: screenshot, capture, visual, playwright
- **Input Schema**: Defined in agent configuration
- **Output Schema**: Defined in agent configuration


### execute_script
- **Description**: Execute custom JavaScript code in the browser context
- **Tags**: javascript, execution, custom, playwright
- **Input Schema**: Defined in agent configuration
- **Output Schema**: Defined in agent configuration


### handle_authentication
- **Description**: Handle various authentication scenarios including basic auth, OAuth, and custom login forms
- **Tags**: authentication, login, security, playwright
- **Input Schema**: Defined in agent configuration
- **Output Schema**: Defined in agent configuration


### wait_for_condition
- **Description**: Wait for specific conditions before proceeding with automation
- **Tags**: wait, synchronization, timing, playwright
- **Input Schema**: Defined in agent configuration
- **Output Schema**: Defined in agent configuration




## Server Configuration

**Port**: 8080

**Debug Mode**: ❌ Disabled



**Authentication**: ❌ Not required


## API Endpoints

The agent exposes the following HTTP endpoints:

- `GET /.well-known/agent.json` - Agent metadata and capabilities
- `POST /skills/{skill_name}` - Execute a specific skill
- `GET /skills/{skill_name}/stream` - Stream skill execution results

## Environment Setup

### Required Environment Variables

Key environment variables you'll need to configure:



- `PORT` - Server port (default: 8080)

### Development Environment


**Flox Environment**: ✅ Configured for reproducible development setup




## Usage

### Starting the Agent

```bash
# Install dependencies
go mod download

# Run the agent
go run main.go

# Or use Task
task run
```


### Communicating with the Agent

The agent implements the A2A protocol and can be communicated with via HTTP requests:

```bash
# Get agent information
curl http://localhost:8080/.well-known/agent.json



# Execute navigate_to_url skill
curl -X POST http://localhost:8080/skills/navigate_to_url \
  -H "Content-Type: application/json" \
  -d '{"input": "your_input_here"}'

# Execute click_element skill
curl -X POST http://localhost:8080/skills/click_element \
  -H "Content-Type: application/json" \
  -d '{"input": "your_input_here"}'

# Execute fill_form skill
curl -X POST http://localhost:8080/skills/fill_form \
  -H "Content-Type: application/json" \
  -d '{"input": "your_input_here"}'

# Execute extract_data skill
curl -X POST http://localhost:8080/skills/extract_data \
  -H "Content-Type: application/json" \
  -d '{"input": "your_input_here"}'

# Execute take_screenshot skill
curl -X POST http://localhost:8080/skills/take_screenshot \
  -H "Content-Type: application/json" \
  -d '{"input": "your_input_here"}'

# Execute execute_script skill
curl -X POST http://localhost:8080/skills/execute_script \
  -H "Content-Type: application/json" \
  -d '{"input": "your_input_here"}'

# Execute handle_authentication skill
curl -X POST http://localhost:8080/skills/handle_authentication \
  -H "Content-Type: application/json" \
  -d '{"input": "your_input_here"}'

# Execute wait_for_condition skill
curl -X POST http://localhost:8080/skills/wait_for_condition \
  -H "Content-Type: application/json" \
  -d '{"input": "your_input_here"}'


```

## Deployment


**Deployment Type**: Manual
- Build and run the agent binary directly
- Use provided Dockerfile for containerized deployment



### Docker Deployment
```bash
# Build image
docker build -t browser-agent .

# Run container
docker run -p 8080:8080 browser-agent
```


## Development

### Project Structure

```
.
├── main.go              # Server entry point
├── skills/              # Business logic skills

│   └── navigate_to_url.go   # Navigate to a specific URL and wait for the page to fully load

│   └── click_element.go   # Click on an element identified by selector, text, or other locator strategies

│   └── fill_form.go   # Fill form fields with provided data, handling various input types

│   └── extract_data.go   # Extract data from the page using selectors and return structured information

│   └── take_screenshot.go   # Capture a screenshot of the current page or specific element

│   └── execute_script.go   # Execute custom JavaScript code in the browser context

│   └── handle_authentication.go   # Handle various authentication scenarios including basic auth, OAuth, and custom login forms

│   └── wait_for_condition.go   # Wait for specific conditions before proceeding with automation

├── .well-known/         # Agent configuration
│   └── agent.json       # Agent metadata
├── go.mod               # Go module definition
└── README.md            # Project documentation
```


### Testing

```bash
# Run tests
task test
go test ./...

# Run with coverage
task test:coverage
```


## Contributing

1. Implement business logic in skill files (replace TODO placeholders)
2. Add comprehensive tests for new functionality
3. Follow the established code patterns and conventions
4. Ensure proper error handling throughout
5. Update documentation as needed

## Agent Metadata

This agent was generated using ADL CLI v0.1.3 with the following configuration:

- **Language**: Go
- **Template**: Minimal A2A Agent
- **ADL Version**: adl.dev/v1

---

For more information about A2A agents and the ADL specification, visit the [ADL CLI documentation](https://github.com/inference-gateway/adl-cli).
