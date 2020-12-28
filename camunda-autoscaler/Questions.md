* How long did it take you to solve the exercise? (Please be honest, we evaluate your answer to this question based on your experience.)

  `about four days 3-4 hours per day to pass both variants and
  one full day to check and describe all these.
  When I started I was confident this task about HPA, I pass it and
  wanted to send, but read task again, I decided you want to see
  kind of operator pod and created one more variant`

* Which additional steps would you take in order to make this code production ready? Why?

  `For default variant`

  `metrics for monitoring,`

  `additional checks for status ready and status started,`

  `move to env variables presets values in the main.go file,`

  `change deployment camunda-deployment readiness probe, a container started and took traffic too early`

  `For hpa variant`

  `move to env variables presets values in the main.go file`

* Which step took most of the time? Why?

  `Find out how works your chalenge task`

  `For default variant debug function getPods with additional checks for a pod status`

  `For HPA variant Read hpa docs for v2beta2 with polices rules`

  `bugs fix`

  `describe all these`
