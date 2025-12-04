buildPlugin(
        failFast: false,
        tests: [skip: true], // tests skipped because no junit reports generated.
        useContainerAgent: false, // Docker builds needed
        configurations: [
                [ platform: "linux", jdk: "21", jenkins: null ],
                [ platform: "windows", jdk: "17", jenkins: null ],
        ])
