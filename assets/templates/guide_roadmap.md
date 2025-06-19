The following summarizes my goal and all critical information gathered so far:

Goal:
{{ .GoalSummary }}

Key details:
{{ .AnswersSummary }}

Instructions:
Build a step-by-step execution roadmap. Every task must be concrete, specific, and actionable. Output as a JSON array—no prose, no comments, no explanations—each object structured as follows:

- "id": Unique integer for this task.
- "description": Clear, actionable task text.
- "project": Create a short, hierarchical project identifier following Taskwarrior best practices, using dot notation (e.g., work.career, personal.health.weight). Derive this from the goal summary—make it concise and specific.
- "tags": Array of relevant tags from {{ .UserTags }}.
- "depends": Array of task ids this task depends on (empty array if none).
- "priority": "High", "Medium", or "Low"—set by urgency or importance.
- "estimate": Estimated duration (e.g., "2h", "3d").
- "resources": List of required resources.
- "risks": Brief summary of risks and mitigation.
- "metrics": What to measure for completion/success.
- "decision_point": true if task is a critical review or decision, false otherwise.

Example format:

[
  {
    "id": 1,
    "description": "Define course outline",
    "project": "work.course",
    "tags": ["pc", "key"],
    "depends": [],
    "priority": "High",
    "estimate": "3d",
    "resources": ["Subject expertise", "Research time"],
    "risks": "Unclear objectives; mitigate by stakeholder review",
    "metrics": "Outline approved by X",
    "decision_point": false
  },
  {
    "id": 2,
    "description": "Record first module",
    "project": "work.course",
    "tags": ["recording"],
    "depends": [1],
    "priority": "Medium",
    "estimate": "2d",
    "resources": ["Recording equipment", "Script"],
    "risks": "Technical issues; mitigate by test recordings",
    "metrics": "Module uploaded and reviewed",
    "decision_point": true
  }
]

Output only the JSON array, nothing else.
