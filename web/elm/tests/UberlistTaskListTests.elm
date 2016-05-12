module UberlistTaskListTests (..) where

import UberlistTask
import UberlistTaskList
import ElmTest exposing (..)
import String
import Html exposing (..)
import Signal
import Array


assertOnListItem : Int -> List a -> (a -> Assertion) -> Assertion
assertOnListItem idx list assCb =
  case (Array.get idx (Array.fromList list)) of
    Nothing ->
      assert False

    Just item ->
      assCb item


all : Test
all =
  let
    model =
      UberlistTaskList.Model
        [ UberlistTask.Model 1 "Hello world" False ]
        0
        Nothing

    modelWithNewTask =
      UberlistTaskList.update UberlistTaskList.CreateNewTask model

    modifiedNewTask =
      UberlistTask.update
        (UberlistTask.UpdateTitleField "New task")
        UberlistTask.newTask

    modelWithModifiedNewTask =
      UberlistTaskList.update
        (UberlistTaskList.ModifyActiveTask
          (UberlistTask.UpdateTitleField "New task")
        )
        modelWithNewTask

    submittedNewTask =
      UberlistTask.update
        (UberlistTask.SubmitNewTitle)
        modifiedNewTask

    modelWithNewTaskSubmitted =
      UberlistTaskList.update
        UberlistTaskList.SubmitNewTask
        modelWithModifiedNewTask

    modelWithNewTaskCancelled =
      UberlistTaskList.update
        UberlistTaskList.CancelNewTask
        modelWithModifiedNewTask

    modelWithTaskSelected =
      UberlistTaskList.update
        (UberlistTaskList.SelectActiveTask 1)
        modelWithNewTaskSubmitted

    modelWithModifiedActiveTask =
      UberlistTaskList.update
        (UberlistTaskList.ModifyActiveTask UberlistTask.EditTitle)
        modelWithTaskSelected
  in
    suite
      "Update UberlistTaskList"
      [ test
          "Can start creating new task"
          (modelWithNewTask.newTask `assertEqual` Just UberlistTask.newTask)
      , test
          "Resets the active task index when creating a new task"
          (modelWithModifiedNewTask.activeTaskIdx `assertEqual` -1)
      , test
          "Can modify the new task"
          (modelWithModifiedNewTask.newTask `assertEqual` Just modifiedNewTask)
      , test
          "Prepends the new task to the list when submitting the new task"
          (assertOnListItem
            0
            modelWithNewTaskSubmitted.tasks
            (\task -> assertEqual task.title modifiedNewTask.title)
          )
      , test
          "Stops modifying the new task when submitting it"
          (assertOnListItem
            0
            modelWithNewTaskSubmitted.tasks
            (\task -> assertEqual task submittedNewTask)
          )
      , test
          "Nullifies the new task field when submitting the new task"
          (modelWithNewTaskSubmitted.newTask `assertEqual` Nothing)
      , test
          "Sets the new taska as active when submitting the new task"
          (modelWithNewTaskSubmitted.activeTaskIdx `assertEqual` 0)
      , test
          "Can cancel the new task"
          (modelWithNewTaskCancelled.newTask `assertEqual` Nothing)
      , test
          "Can select a task"
          (modelWithTaskSelected.activeTaskIdx `assertEqual` 1)
      , test
          "Cannot select a task when creating a new task"
          ((UberlistTaskList.update
              (UberlistTaskList.SelectActiveTask 0)
              modelWithModifiedNewTask
           )
            `assertEqual` modelWithModifiedNewTask
          )
      , test
          "Can modify the selected task"
          (assertOnListItem
            1
            modelWithModifiedActiveTask.tasks
            (\task -> assert task.titleIsField)
          )
      ]
