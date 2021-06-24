---
name: 🐞 Bug report
about: Report a bug
title: ''
labels: ''
assignees: ''

---

/kind bug

<!--

Welcome! - We kindly ask you to:

  1. Fill out the issue template below 
  2. Use Gitter if you have a question rather than a bug or feature request.
    Gitter community channel is at: https://gitter.im/devfile/community

Thanks for understanding, and for contributing to the project!

-->

### Which area this bug is related to?
<!--
    Uncomment appropriate `/area` lines, and delete the rest.
    For example, `> /area api` would simply become: `/area api`
-->

> /area ci
> /area api
> /area library
> /area registry
> /area devworkspace
> /area integration-tests
> /area test-automation


## What versions of software are you using?

**Operating System:**

**Go Pkg Version:**


## Bug Summary
**Describe the bug:**
<!-- 
    A clear and concise description of what the bug is.
-->

**To Reproduce:**
<!--
    Steps to reproduce the behavior.
-->


## Expected behavior
<!--
    A clear and concise description of what you expected to happen.
-->

## Any logs, error output, screenshots etc? Provide the devfile that sees this bug, if applicable.
<!--
To get logs:
    ci: please copy the github workflow output
    api: please copy the terminal output
    library: please copy the terminal output
    registry: follow instruction under "Collecting Logs" to find log: https://github.com/devfile/registry-support/blob/main/TROUBLESHOOTING.md
    devworkspace: copy the logs from the controller (kubectl logs deploy/devworkspace-controller -n $NAMESPACE)
    test-automation: 
        api: follow instruction under "Running tests locally" to find test log: https://github.com/devfile/api/tree/main/test
        library: follow instruction under "Running the tests locally" to find test log: https://github.com/devfile/library/tree/main/tests
        devworkspace: copy the logs from the controller (kubectl logs deploy/devworkspace-controller -n $NAMESPACE)
    integration-tests: please copy the build log under prow ci result for QE ingetration tests
-->


## Additional context
<!--
    Add any other context about the problem here.
-->
**Any workaround?**
<!--
    Describe the workaround if applicable.
-->

**Suggestion on how to fix the bug**
<!--
    Provide suggestion on how to fix the bug upon your investigation, if applicable.
-->


