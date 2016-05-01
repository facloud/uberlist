module Main (..) where

import Signal exposing (Signal)
import ElmTest exposing (consoleRunner, suite)
import Console exposing (IO, run)
import Task
import UberlistTaskTests
import UberlistTaskListTests


console : IO ()
console =
  consoleRunner
    (suite
      "Uberlist tests"
      [ UberlistTaskTests.all
      , UberlistTaskListTests.all
      ]
    )


port runner : Signal (Task.Task x ())
port runner =
  run console
