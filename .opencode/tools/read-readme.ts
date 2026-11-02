import { tool } from "@opencode-ai/plugin"
import { readFile } from "node:fs/promises"
import path from "node:path"

export default tool({
  description:
    "Reads the first rows of the project README.md from the worktree root",
  args: {
    lines: tool.schema.number()
      .int()
      .min(1)
      .max(50)
      .default(10)
      .describe("Number of rows to return"),
  },
  async execute(args, context) {
    const readmePath = path.join(context.worktree, "README.md")
    const content = await readFile(readmePath, "utf8").catch(() => null)
    if (content === null) {
      return `README.md not found at ${readmePath}`
    }
    const rows = content.split("\n").slice(0, args.lines)
    return {
      title: `README.md: first ${args.lines} rows`,
      output: rows.join("\n"),
    }
  },
})
