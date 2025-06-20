{{template "user_context.md" .}}

# Provided Data

## Tasks
{{ range $index, $task := .Tasks }}
### Task {{ add $index 1 }}
- Description: {{ $task.Description }}
- Tags: [{{ range $task.Tags }}{{ . }},{{ end }}]
- Project: {{ if $task.Project }}{{ $task.Project }}{{ else }}(none){{ end }}
- Due: {{ if $task.DueDate }}{{ $task.DueDate }}{{ else }}(none){{ end }}

{{ end }}

---

# Objectives:

For each task, provide analysis following these guidelines:

## 1. Analyze
 Only apply tags presented to you: Suggest up to 5 of them. Tags have a + as Prefix. Dont suggest tags you are not sure of adding. Keep the project the same. If no project is set, assign the best one of those that are presented to you or create one using dot notation (`personal.health`, `wrk.career`, etc.).

## 2. Refine
Start task with a verb. Keep it short and clear. Only extend if it's vague. Make it actionable and concrete.

## 3. Contextual Info (Optional)
Return these only if clearly relevant. Max one sentence each:
{{ range .UserContext.UserAnnotations }}
- {{ .Name }}: {{ .Description }} {{ end }}

## 4. Subtask
If task is broad or complex enough, split into 3–5 actionable subtasks (same refinement rules). Don't base tags/goal alignment on subtasks.

---

# Output (JSON)

{
  "task_analyses": [
    {
      "task_index": 1,
      "suggested_tags": ["+tag1", "+tag2"],
      "goal_alignment": "How this supports user goals.",
      "project": "project.name",
      "refined_task": "Refined task text here",
      "additional_infos": {{ .ExampleOutput }},
      "subtasks": [
        "Do this",
        "Then this",
        "Finally this"
      ]
    }
  ]
}

Only answer with valid json. Dont use projects as tags or tags as projects. Provide analysis for all {{ len .Tasks }} tasks.