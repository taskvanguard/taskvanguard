<a id="readme-top"></a>

<!-- <div align="center"> -->
<!-- [![Contributors][contributors-shield]][contributors-url] -->
<!-- [![Forks][forks-shield]][forks-url] -->
<!-- [![Stargazers][stars-shield]][stars-url] -->
<!-- [![Issues][issues-shield]][issues-url] -->
<!-- [![Task][license-shield]][license-url] -->
<!-- [![LinkedIn][linkedin-shield]][linkedin-url] -->
<!-- </div> -->

<!-- PROJECT LOGO -->
<div align="center">
  <a href="https://github.com/taskvanguard/taskvanguard">
    <img src="docs/images/logo.png" alt="Logo" width="200" height="200">
  </a>

<h3 align="center">Task Vanguard [BETA 0.2.1]</h3>

  <p align="center">
Your tactical advisor at the frontlines of your goals - powered by AI using TaskWarrior.
    <br>

**TaskVanguard** is a lightweight, fast, highly configurable CLI wrapper for [TaskWarrior](https://taskwarrior.org/), written in Go. It brings AI-powered suggestions, smart tagging, goal management and cognitive support using any OpenAI-compatible LLM API.

<br>
<a href="https://buymeacoffee.com/taskvanguard">Donate</a>
&middot;
<a href="https://github.com/taskvanguard/taskvanguard/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
&middot;
<a href="https://github.com/taskvanguard/taskvanguard/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>

<!-- <div align="center"> -->
<!-- [![Contributors][contributors-shield]][contributors-url] -->
<!-- [![Forks][forks-shield]][forks-url] -->
[![Go][Go-shield]][Go.dev]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
<!-- [![Task][license-shield]][license-url] -->
<!-- [![LinkedIn][linkedin-shield]][linkedin-url] -->
<!-- </div> -->

</div>


<!-- ABOUT THE PROJECT -->
## What is TaskVanguard?

Use `vanguard add <task>` just like `taskwarrior add <task>`. TaskVanguard creates the task, then suggests improvements using an LLM (OpenAI, Deepseek, etc).

<div align="center">

![Product Name Screen Shot][product-screenshot]

</div>

### Features

✨ **Add:** AI-Enhanced Task Creation: Improves task titles, tags, project and annotations for clarity and relevance.<br>
🎯 **Spot:** Do the Right Thing Next: Identifies the most impactful next task. Based on urgency, context, mood etc.<br>
🧭 **Guidance:** Generate concrete, step-by-step roadmaps to achieve goals using LLM-backed planning.<br>
⛰️ **Goal Management:** Link tasks to long-term objectives and maintain alignment with your broader mission.<br>
📦 **Batch Analysis:** Refactor, annotate entire task backlogs by tags or projects at once.<br>
🗡️ **Subtask Splitting:** Suggests splitting up vague tasks and suggests clear, actionable subtasks.

**Tip:** You can stop certain tasks from being sent to the LLM by blacklisting tags or projects.   


## Why TaskVanguard?

- **Stalled by stale high-priority tasks?** Reframe what moves your mission forward.
- **Tasks too broad or unclear?** Break them into precise, executable steps.
- **Spending time on structure instead of action?** Let the system handle the overhead.
- **Unsure what’s worth doing now?** Surface the tasks with real leverage.

**⚔️ TaskVanguard** fills those gaps using LLMs for real cognitive support. It’s especially useful for ADHD-driven procrastination: it reduces friction to start and helps reframe daunting tasks.


<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Quickstart


### 1. Prerequisites

- [TaskWarrior](https://taskwarrior.org/) installed and initialized.
- Go 1.21+
- API key for your preferred LLM (OpenAI, OpenRouter, Deepseek etc).


### 2. Install

```sh
go install github.com/taskvanguard/taskvanguard/cmd/vanguard@latest
```

or

```sh
git clone https://github.com/taskvanguard/taskvanguard.git
cd taskvanguard
go build ./cmd/vanguard
```

<!-- Or download the latest release from [Releases](https://github.com/taskvanguard/taskvanguard/releases). -->

### 3. Configure

```sh
./vanguard init
```

- Creates a default config at ~/.config/taskvanguard/vanguardrc.yaml
- Prompts for your LLM API key
- Suggests shell alias: `alias tvg="vanguard"`
- Ensures TaskWarrior is set up with task command working

### 4. Usage

```sh
# Add a task (AI-augmented)
vanguard add "refactor onboarding flow" project:work

# Analyze all tasks (offers refactoring & suggests annotation)
vanguard analyze

# Get the one highly important task to do next (spotlight)
vanguard spot
```

See ``vanguard --help`` for full options

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Usage

**Tags**

The following Tags are added by default:

- `+sb` **Snowballing**:  Task has the potential for compounding effects (positive or negative); tackling it may unlock cascading benefits or risks.
- `+cut` **Cut**: Task can save time or money in the future.
- `+fast` **Fast**: Task can be finished quickly and requires minimal setup.
- `+key` **Key**: High-impact task that directly drives one or more of your goals.


**Annotations**

By default the following annotations are added via LLM to the tasks you refactor with TaskVanguard:

- `short_reward:` immediate benefit
- `long_reward:` strategic benefit
- `risk:` what if not done
- `tip:` practical, actionable, insightful


**Goal Management:**

- Add  major goals as tasks in `project:goals` just like you would create any other task.
- Link any task to a goal for automatic relationship tracking by using an uda (`vanguard goals link <task_id> <goal_id>`).
- Use the `vanguard goals` command for comprehensive goal management.


## Commands

| Command            | Description                                     |
| ------------------ | ----------------------------------------------- |
| `vanguard init`    | Configure API, shell aliases, default settings  |
| `vanguard add`     | Creates tasks with AI-augmented improvements                      |
| `vanguard analyze` | Provides LLM-driven review, tags, and refactoring |
| `vanguard spot`    | Surfaces the single best task to do next        |
| `vanguard goals`   | Manage goals and link tasks to achieve them     |


### Init

- Creates Default Config
- Suggests adding aliases for taskVanguard to .bashrc
- Asks for API Key
- Backups Tasks & Config
- Configures tags and annotations used by TaskVanguard
- Adds urgency coefficient factors for +sb, +cut, +fast, +key
- Adds urgency coefficient factors for goals (heavy minus)

### Add

Processes a task add command with an LLM to suggest better title, tags and annotations or splitting into subtasks.

- `--no-tags` disables LLM suggestioning tags. (config on/off)
- `--no-subtasks` disables subtask splitting. (config on/off)
- `--no-annotations` disables LLM suggestioning annotatians. (config on/off)

### Spot

Takes into account mood and context. You can use `--no-prompt` and use it for system notifications. If you type in a mood and context it will remember that for 4 hours.

- `--no-prompt` skips questions about context and mood 
- `--mood <mood>` provide the mood you are in (for example energetic)
- `--context <context>` provide the context (for example location) 
- `--refresh` overwrite the cache with the context and mood you provide


### Analyze

Analyzes either a specific task or a list of tasks and suggests improved task descriptions and tag assignments. If you analyze a specific task it suggests annotations and linking to a specific goal.

- `--batch-editor` opens your $EDITOR, allowing to edit all the task suggestions at once before applying them by saving and quits
- `--interactive` apply suggestions one by one for each task

### Goals

Goals are primarily managed in the background. When you use `vanguard guide`, a goal is defined and a step-by-step roadmap is generated to help you achieve it. All related tasks are automatically linked to that goal. By associating tasks with goals, TaskVanguard can better understand the context in which each task exists-going beyond simple tagging (like +sb or +key). Goals are actually regular tasks within a special project (named goals by default, but customizable in your config).

#### Goal Commands

| Command                    | Description                                    |
| -------------------------- | ---------------------------------------------- |
| `vanguard goals list`      | List all your goals                           |
| `vanguard goals add <desc>` | Create a new goal                             |
| `vanguard goals show <id>` | Show detailed information about a goal/task   |
| `vanguard goals modify <id> <args>` | Modify an existing goal               |
| `vanguard goals delete <id>` | Delete a goal                               |
| `vanguard goals link <id1> <id2>` | Link a task to a goal (order-agnostic) |
| `vanguard goals unlink <id1> <id2>` | Remove task-goal link               |
| `vanguard goals links <id>` | Show all tasks linked to a goal or goal linked to a task |

#### Goal Usage

```bash
# Create a new goal
vanguard goals add "complete certification in cloud architecture"

# Link an existing task to a goal (works both ways)
vanguard goals link <task_id> <goal_id>   # task 123 -> goal 456
vanguard goals link <goal_id> <task_id>   # same result

# See what tasks are linked to a goal
vanguard goals links <goal_id>

# See which goal a task is linked to
vanguard goals links <task_id>

# Modify a goal like any TaskWarrior task
vanguard goals modify <goal_id> "priority:H due:2024-12-31"
```

#### Goal Features

- **TaskWarrior Integration**: Goals are stored as regular TaskWarrior tasks in a dedicated project
- **Flexible Linking**: Link any task to any goal using the `goal` UDA (User Defined Attribute)
- **Order-Agnostic Commands**: Link/unlink commands work regardless of argument order
- **Relationship Tracking**: Easily see which tasks contribute to which goals and vice versa
- **Configurable Project**: Set your own goal project name via `goal_project_name` in config
- **Full TaskWarrior Compatibility**: Goals support all TaskWarrior features (tags, priority, due dates, etc.)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- 
## CronJob

You can set up a cronjob like this:
``0 */12 * * * vanguard spot --no-prompt`` -->

## Alias

```bash
alias tvg="vanguard"
alias tvga="vanguard add"  # for completions
```

## Config

`~/.config/taskvanguard/vanguardrc.yaml` Key settings:
- `debug`: Enable verbose logging for troubleshooting.
- `enable_llm`: Pretty much necessary.
- `split_tasks`: Allow LLM to suggest subtask splits.
- `auto_import_tags`: Retrieve the tags that are in use for suggestions.
- `auto_import_projects`: retrieve the projects that are in use for suggestions.
- `enable_lowercase`: Make all suggestions lowercase.
- `enable_tagging`: Enables LLM tagging suggestions. 
- `enable_annotations`: Enables annotation suggestions.
- `enable_goals`: Enables LLM linking tasks to projects.
- `goal_project_name`: Name of your goals project.
- `task_processing_batch_size`: Number of tasks processed at once (default: 15).
- `task_import_limit`: Max tasks to import for analysis (default: 999).

```yaml
settings:
    debug: true
    enable_llm: false
    split_tasks: true
    auto_import_tags: true
    auto_import_projects: true
    enable_lowercase: true
    enable_tagging: true
    enable_annotations: true
    enable_goals: true
    goal_project_name: "goals"
    task_import_limit: 999
    task_processing_batch_size: 15
llm:
    provider: openai
    api_key: "<YOUR_API_KEY_HERE>"
    model: gpt-3.5-turbo
    base_url: https://openrouter.ai/api/v1
filters:
    tag_filter_mode: "blacklist"
    tag_filter_tags: ["private", "confidential"]
    project_filter_mode: "blacklist"
    project_filter_projects: ["pers.secret", "work.secret"]
tags:
    cut:
        desc: Task has the potential to save time or cost in the future
        urgency_factor: 1.2
    key:
        desc: Task is impacting goals of the user
        urgency_factor: 1.2
    fast:
        desc: Task is probably done in very short time (10 Mins or less)
        urgency_factor: 1.2
    sb:
        desc: Task is potentially snowballing positively or negatively and offers high roi
        urgency_factor: 1.3
annotations:
    short_reward:
        label: Short Reward
        symbol: ●
        description: immediate benefit
    long_reward:
        label: Long Reward
        symbol: ●
        description: strategic benefit
    risk:
        label: Risk
        symbol: ●
        description: if not done
    tip:
        label: Tip
        symbol: ●
        description: practical, actionable, insightful

```

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- CONTRIBUTING -->
## Contributing

By submitting code or other contributions to this project, you agree to the Contributor License Agreement (CLA). See [CLA.md](./CLA.md) for details.

PRs and issues welcome. Open an enhancement issue or fork and create a pull request.

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- LICENSE -->
## License

This project is licensed under the GNU Affero General Public License v3.0.
See `LICENSE.txt`.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

xarc - [@xarcdev](https://x.com/xarcdev) - taskvanguard@xarc.dev

Project Link: [https://github.com/taskvanguard/taskvanguard](https://github.com/taskvanguard/taskvanguard)

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/taskvanguard/taskvanguard.svg?style=for-the-badge
[contributors-url]: https://github.com/taskvanguard/taskvanguard/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/taskvanguard/taskvanguard.svg?style=for-the-badge
[forks-url]: https://github.com/taskvanguard/taskvanguard/network/members
[stars-shield]: https://img.shields.io/github/stars/taskvanguard/taskvanguard.svg?style=for-the-badge
[stars-url]: https://github.com/taskvanguard/taskvanguard/stargazers
[issues-shield]: https://img.shields.io/github/issues/taskvanguard/taskvanguard.svg?style=for-the-badge
[issues-url]: https://github.com/taskvanguard/taskvanguard/issues
[license-shield]: https://img.shields.io/github/license/taskvanguard/taskvanguard.svg?style=for-the-badge
[license-url]: https://github.com/taskvanguard/taskvanguard/blob/master/LICENSE.txt
<!-- [linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555 -->
<!-- [linkedin-url]: https://linkedin.com/in/linkedin_username -->
[product-screenshot]: docs/images/screenshot.webp

[Go.dev]: https://go.dev/
[Go-shield]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white



