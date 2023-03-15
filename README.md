# gox

`gox` is a package for executing multiple functions concurrently in Go, with optional error handling. It provides a simple interface to start and wait for the completion of one or more functions, and collect any errors that may occur during their execution.

## Installation

To use gox, you need to have Go installed and set up on your system. Then you can install the package with the following command:

```sh
go get github.com/afixer/gox
```

## Usage

### Simple concurrent execution

To start executing one or more functions concurrently, you can use the Run function of the gox package:

```go
import (
 "fmt"
 "github.com/afixer/gox"
)

func main() {
 gox.Run(
  func() { fmt.Println("Hello") },
  func() { fmt.Println("World") },
 ).Wait()
 fmt.Println("Done")
}
```

This will start two goroutines, one for each function, and print "Hello" and "World" in any order. The Run function returns a *gox.goX struct that you can use to wait for the completion of all functions with the Wait method. In this example, we chain the Wait method to the Run method call to wait for the completion of all functions before exiting.

### Concurrent execution with error handling

If you need to execute functions that may return errors, you can use the RunE function of the gox package:

```go
package main

import (
 "errors"
 "fmt"
 "github.com/afixer/gox"
)

func main() {
 errMap := gox.RunE("foo", func() error {
   return errors.New("foo error")
  }).RunE("bar", func() error {
   return errors.New("bar error")
  }).WaitE()

 fmt.Println(errMap)
}
```

This will start two goroutines, one for each function, and collect any errors that occur during their execution. The RunE method takes a function that returns an error, and a string that represents its name. The returned *gox.goE struct provides a WaitE method that returns a map of errors, keyed by their function name, or nil if no errors occurred.
