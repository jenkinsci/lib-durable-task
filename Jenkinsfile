#!/usr/bin/env groovy

/* `buildPlugin` step provided by: https://github.com/jenkins-infra/pipeline-library */
// tests skipped because no junit reports generated.
buildPlugin(failFast: false, tests: [skip: true], configurations: [
        [ platform: "docker", jdk: "11", jenkins: null ],
        [ platform: "docker-windows", jdk: "11", jenkins: null ],
])
