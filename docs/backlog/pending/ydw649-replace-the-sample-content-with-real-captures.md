---
title: 'Replace the sample content with real captures'
priority: p2
section: landing
topics:
    - content
    - accuracy
---
The page labels three things as sample: call 3's reply and the post-install files (shapes taken from `internal/mcp/install_tool.go`) and the approval dialog (drawn, not a capture of the client). Replacing them needs a real `install_skills` run and `rsk install` output, which requires installs, and the agent does not install without explicit approval.

Done when each is a capture, labelled `real run`, with the `facts.ts` constants and the design doc's real-versus-sample table updated.
