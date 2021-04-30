#!/usr/bin/env groovy

/* `buildPlugin` step provided by: https://github.com/jenkins-infra/pipeline-library */
// tests skipped on windows because no surefire reports generated. However, unit tests are still being run and will
// fail build if failed
buildPlugin(failFast: false, configurations: [
        [ platform: "docker", jdk: "8", jenkins: null ],
        [ platform: "windock", jdk: "8", jenkins: null, skipTests: true ]
])