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

<h3 align="center">Task Vanguard</h3>

  <p align="center">
Your tactical advisor at the frontlines of your goals - powered by AI using TaskWarrior.
    <br>

**TaskVanguard** is a lightweight, fast, highly configurable CLI wrapper for [TaskWarrior](https://taskwarrior.org/), written in Go. It brings AI-powered suggestions, smarter tagging, goal alignment, and real cognitive support to your daily workflow using any OpenAI-compatible LLM API.

<br>
<a href="https://github.com/taskvanguard/taskvanguard">Donate</a>
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

Use taskvanguard add <task> like taskwarrior add <task>. TaskVanguard creates the task, then suggests improvements using your LLM of choice (OpenAI, Deepseek, etc).

<div align="center">

![Product Name Screen Shot][product-screenshot]

</div>

- **AI-enhanced task creation:** Auto-improves task titles, tags, project assignment, and annotations.
- **Context-aware prioritization:** Surfaces high-impact tasks and explains *why* you should do them now.
- **Goal linking:** Assign tasks to goals and keep your focus on long-term outcomes.
- **Batch analysis:** Mass-edit, annotate, or refactor your existing backlog using your LLM of choice.
- **Subtask splitting:** Detects “too big” tasks and offers splits.
- **Highly configurable:** Full YAML config, tag/project/goal blacklists, and fine-grained toggles.

**Tip:** You can stop the LLM from processing certain Tasks by blacklisting tags or projects, configure what annotations it should generate for your tasks, if any.   


## Why TaskVanguard?

- **Tired of TaskWarrior leaving you stuck with “stale” high-priority tasks at the top of your list?**
- **Want to break big tasks into actionables without manual splitting?**
- **Need to shorten the amount of time spent thinking about tags and actionable task description?**
- **Wish you had a second opinion on which tasks are high impact and for motivation?**

**TaskVanguard** fills those gaps using LLMs for real cognitive support. It’s especially useful for ADHD-driven procrastination: it reduces friction to start and helps reframe daunting tasks.


<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Quickstart


### 1. Prerequisites

- [TaskWarrior](https://taskwarrior.org/) installed and initialized.
- Go 1.21+
- API key for your preferred LLM (OpenAI, OpenRouter, Deepseek etc).


### 2. Install

```sh
git clone https://github.com/taskvanguard/taskvanguard.git
cd taskvanguard
go build -o vanguard
```

<!-- Or download the latest release from [Releases](https://github.com/taskvanguard/taskvanguard/releases). -->

### 3. Configure

```sh
./vanguard init
```

- Creates a default config at ~/.config/taskvanguard/vanguardrc.yaml
- Prompts for your LLM API key
- Suggests shell alias: `alias tvg="vanguard"`
- Ensure TaskWarrior is set up with task command working.

### 4. Usage

```sh
# Add a task (AI-augmented)
vanguard add "Refactor onboarding flow" +project:infra +sb

# Analyze all tasks (refactoring, annotation suggestions)
vanguard analyze --editor

# Get the one thing to do next (spotlight)
vanguard spot
```

See ``vanguard --help`` for full options

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Commands

| Command            | Description                                     |
| ------------------ | ----------------------------------------------- |
| `vanguard init`    | Configure API, shell aliases, default settings  |
| `vanguard add`     | Creates tasks with AI-augmented improvements                      |
| `vanguard analyze` | Provides LLM-driven review, tags, and refactoring |
| `vanguard spot`    | Surfaces the single best task to do next        |

**TaskVanguard introduces several high-leverage tags by default:**

- `+sb` **Snowballing**:  Task has the potential for compounding effects (positive or negative); tackling it may unlock cascading benefits or risks.
- `+cut` **Cut**: Task can save time or money in the future.
- `+fast` **Fast**: Task can be finished quickly and requires minimal setup.
- `+key` **Key**: High-impact task that directly drives one or more of your goals.

**Goal Management:**

- Create major goals as tasks in `project:goals` just like you would add any other task.
- Link any task to a goal for automatic relationship tracking.


### Init

- Create Default Config
- Suggesting adding aliases for taskVanguard to .bashrc
- Adding API Key
- Backup Tasks & Config
- Setup tags used by task vanguard
- Add urgency coefficient factors for +sb, +cut, +fast, +key
- Add urgency coefficient factors for goals (heavy minus)
- Suggests to run `analyze`

### Add

Processes the task with an LLM to suggest a better title, tags, and subtasks.

- `--no-tags` disables LLM tag suggestions.
- `--no-subtasks` disables subtask splitting. Use config to disable task splitting by default.

### Spot

Takes into account mood and context. You can use --no-prompt and use it in notifications. If you type in a mood and context it will remember that for 4 hours.

- `--no-prompt` skips questions about context and mood 
- `--mood` provide the mood you are in
- `--context` provide the context (for example location) 
- `--refresh` overwrite the cache with the context and mood you provide


### Analyze

Analyzes either a specific or all your tasks and presents improved task description and tag assignment. If you analyze a specific task it suggests annotations and linking to a specific goal.

- `--batch-editor` opens your $EDITOR, allowing to edit all the task suggestions at once before applying them by saving and quits
- `--interactive` apply suggestions one by one for each task

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

⚠️ By submitting code or other contributions to this project, you agree to the Contributor License Agreement (CLA). See [CLA.md](./CLA.md) for details.

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



