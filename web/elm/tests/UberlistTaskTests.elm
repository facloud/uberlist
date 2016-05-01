module UberlistTaskTests (..) where

import UberlistTask
import ElmTest exposing (..)
import String
import Html exposing (..)
import Signal


all : Test
all =
  let
    defaultModel =
      UberlistTask.Model 1 "Hello world" False

    editableModel =
      UberlistTask.update UberlistTask.EditTitle defaultModel

    modifiedModel =
      UberlistTask.update (UberlistTask.UpdateTitleField "New title") editableModel

    newModel =
      UberlistTask.update UberlistTask.SubmitNewTitle modifiedModel
  in
    suite
      "Update UberlistTask"
      [ test
          "Can start editting"
          (assert editableModel.titleIsField)
      , test
          "Can modify the title field"
          (modifiedModel.title `assertEqual` "New title")
      , test
          "Can save the title filed"
          (newModel.titleIsField `assertEqual` False)
      ]
