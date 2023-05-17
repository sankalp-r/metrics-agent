# metrics-agent

## How to run and test

* To update dependencies, run `make deps-update` from Makefile.


* To run the unit-tests, run `make test` from Makefile.


* To run the agent locally, run `make run` from Makefile.

  `config` file must be passed as an argument `-config=<path-to-config-file >` (refer Makefile).

  A sample `config.yaml` is present in the root-directory for reference.

```yaml
  # specify all the http metrics source
  httpsources:
   - endpoints: http://localhost:8080/metrics
     headers:
      x-api-key : abcd
   - endpoints: http://localhost:8080/metricsv2/status
     headers:
      x-api-key : abcd

  # sample frequency is in seconds
  sampleFrequency: 5

  # output file name
  targetOutputFile: output_file
```