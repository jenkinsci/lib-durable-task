#!/usr/bin/env groovy

/* `buildPlugin` step provided by: https://github.com/jenkins-infra/pipeline-library */
// tests skipped because no junit reports generated.
buildPlugin(
        failFast: false,
        tests: [skip: true],
        useContainerAgent: false, // Set to `false` if you need to use Docker for containerized tests
        configurations: [
                [ platform: "linux", jdk: "21", jenkins: null ],
                [ platform: "windows", jdk: "17", jenkins: null ],
        ])
