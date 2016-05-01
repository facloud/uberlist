module UberlistTaskList (..) where

import UberlistTask
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (on, onClick)
import Array


type alias Model =
  { tasks : List UberlistTask.Model
  , activeTaskIdx : Int
  , newTask : Maybe UberlistTask.Model
  }


getActiveTask : Model -> Maybe UberlistTask.Model
getActiveTask model =
  case model.newTask of
    Just newTask ->
      Just newTask

    Nothing ->
      Array.get model.activeTaskIdx (Array.fromList model.tasks)


type Action
  = CreateNewTask
  | SubmitNewTask
  | CancelNewTask
  | SelectActiveTask Int
  | ModifyActiveTask UberlistTask.Action


update : Action -> Model -> Model
update action model =
  case action of
    CreateNewTask ->
      { model
        | newTask = Just (UberlistTask.newTask)
        , activeTaskIdx = -1
      }

    SubmitNewTask ->
      case model.newTask of
        Nothing ->
          model

        Just newTask ->
          { model
            | tasks =
                List.append [ newTask ] model.tasks
            , newTask = Nothing
            , activeTaskIdx = 0
          }

    CancelNewTask ->
      { model
        | newTask = Nothing
      }

    SelectActiveTask idx ->
      case model.newTask of
        Nothing ->
          { model | activeTaskIdx = idx }

        Just _ ->
          model

    ModifyActiveTask taskAction ->
      case model.newTask of
        Nothing ->
          let
            arr =
              Array.fromList model.tasks

            maybeActiveTask =
              getActiveTask model
          in
            case maybeActiveTask of
              Nothing ->
                model

              Just activeTask ->
                { model
                  | tasks =
                      Array.toList
                        (Array.set
                          model.activeTaskIdx
                          (UberlistTask.update taskAction activeTask)
                          arr
                        )
                }

        Just currentNewTask ->
          { model
            | newTask =
                Just (UberlistTask.update taskAction currentNewTask)
          }


view : Signal.Address Action -> Model -> Html
view address model =
  div
    [ class "container" ]
    [ div
        [ id "TaskSelector", class "col-lg-5 list-group" ]
        (List.indexedMap (viewListEntry address model.activeTaskIdx) model.tasks)
    , viewActiveTask address model
    ]


viewActiveTask : Signal.Address Action -> Model -> Html
viewActiveTask address model =
  case (getActiveTask model) of
    Nothing ->
      div [] []

    Just activeTask ->
      UberlistTask.view (Signal.forwardTo address ModifyActiveTask) activeTask


viewListEntry : Signal.Address Action -> Int -> Int -> UberlistTask.Model -> Html
viewListEntry address activeIdx idx task =
  let
    activeClass =
      if activeIdx == idx then
        " active"
      else
        ""
  in
    a
      [ href "#"
      , class ("list-group-item" ++ activeClass)
      , onClick address (SelectActiveTask idx)
      ]
      [ p [ class "list-group-item-text" ] [ text task.title ]
      ]
