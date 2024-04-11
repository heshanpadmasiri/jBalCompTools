# Build jBallerina compiler
+ [x] Optionally set build flags
+ [ ] Make default flags configurable from the run control file

# Run Ballerina source
+ [x] Run projects
+ [x] Run individual files
+ [x] Remote debug runtime
    + [ ] Make default port part for run control file

# Just compile Ballerina source
+ [ ] Compile projects
+ [ ] Compile individual files
+ [ ] Remote debug compiler

# CI
+ [ ] Create unit tests to validate each command
+ [ ] Setup Github workflow to run them
+ [ ] Create Github workflow to build the tool for different operating systems and create release
    + Automatically create a release if the CI is passing

# Disassemble generated jar file
+ [ ] Extend underlying compile command to then disassemble the generated jar file
+ [ ] Given the method and class name show the bytecode

# Benchmark
## Direct measurements
+ [ ] Measure compile time
+ [ ] Measure execution time

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