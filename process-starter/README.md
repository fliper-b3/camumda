# Camunda Process Starter

The Camunda Process Starter is a simple application that creates process
instances on Camunda BPM. It will start random processes based on
a schedule defined by the user. By default, it starts five processes every
two seconds.

There are different environment variables that can influence its
behaviour:

 - `CAMUNDA_HOST` (default: `http://localhost:8080/engine-rest/`) - URL to
   Camunda's engine.
 - `N_PROCESS_STARTED` (default: 5) - Number of processes to be started
 - `QUIET_TIME_S` (default: 2) - Wait time before starting new processes
