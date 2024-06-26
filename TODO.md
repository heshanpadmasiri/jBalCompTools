# Build jBallerina compiler
+ [x] Optionally set build flags
+ [ ] Make default flags configurable from the run control file
+ [ ] Make it possible to change the version from the tool instead of having to change it in `build.gradle`

# Run Ballerina source
+ [x] Run projects
+ [x] Run individual files
+ [x] Remote debug runtime
    + [ ] Make default port part for run control file
    + [x] Show the output while the command is running

# Just compile Ballerina source
+ [x] Compile projects
+ [x] Compile individual files
+ [x] Remote debug compiler

# Run unit tests
+ [x] Run unit tests on the given project

# CI
+ [ ] Create unit tests to validate each command
+ [x] Setup Github workflow to run them
+ [x] Create Github workflow to build the tool for different operating systems and create release
    + [ ] Automatically create a release if the CI is passing

# Disassemble generated jar file
+ [x] Extend underlying compile command to then disassemble the generated jar file
+ [ ] Given the method and class name show the bytecode

# Benchmark
## Direct measurements
+ [x] Measure compile time
+ [x] Measure execution time
+ [ ] Add ability to configure benchmarks (number of times to run, warmup runs etc.)

## Compare performance
+ [ ] Compare against a given ballerina release version
+ [ ] Compare against a given "pack"

# Perf
+ [ ] Given the method and class name show optimizations
    + [ ] Show when each optimizing compiler got triggered on that method
    + [ ] Show the optimized assembly generated for that method

# Native helper
+ [ ] Given ballerina source method name find the java method name
    + [ ] Handle large method splitter
+ [ ] Show ballerina code and bytecode side by side
