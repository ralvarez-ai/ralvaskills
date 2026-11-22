---
title: 'The git runner honors url.*.insteadOf, so the URL shown may not be the URL used'
priority: p2
section: ralvaskills
topics:
    - rsk
    - sources
related:
    - ralvaskills#6w4kzn
---
Found in the SSH acceptance run. With a failing `ssh` first on `PATH`, `rsk catalog` failed on `@anthropics/skills`, whose URL is `https://github.com/anthropics/skills`: the owner's global git config has `url.git@github.com:.insteadOf https://github.com/`, and the hardened runner does not ignore it. So the HTTPS fallback in the implicit-source prompt, the official source and any HTTPS source travel over SSH on that machine, and a rule such as `url.<host-alias>:planttalk/.insteadOf https://github.com/planttalk/` can send a granted URL to a different host. The trust grant and the gate are keyed on the URL rsk holds, before git rewrites it.

This is the owner's own configuration, which the design accepts (trust rule 1: what git sends is the user's git configuration), but the prompt claims a transport it may not use. Decide: show the effective URL (`git ls-remote --get-url`) in the implicit prompt and in `rsk source test`, or run the runner with `GIT_CONFIG_GLOBAL` and `GIT_CONFIG_SYSTEM` unset, which would also drop the owner's credential helpers and so needs its own decision.

Done when the URL shown to the user and the host contacted agree, or the difference is documented where the prompt is.
