module Main (..) where

import UberlistTask
import UberlistTaskList
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick)
import Signal
import StartApp.Simple as StartApp


type alias Model =
  UberlistTaskList.Model


type alias Action =
  UberlistTaskList.Action


view : Signal.Address Action -> Model -> Html
view address model =
  div
    []
    [ nav
        [ class "navbar navbar-default" ]
        [ div
            [ class "container" ]
            [ div
                [ class "navbar-header" ]
                [ a
                    [ class "navbar-brand", href "#" ]
                    [ div [ id "Logo" ] [ h1 [] [ text "UberList" ] ]
                    ]
                ]
            , viewNewStoryButtons address model
            ]
        ]
    , UberlistTaskList.view address model
    ]


viewNewStoryButtons address model =
  div
    [ class "navbar-form navbar-right" ]
    (case model.newTask of
      Nothing ->
        [ button
            [ class "btn btn-success"
            , onClick address UberlistTaskList.CreateNewTask
            ]
            [ span [ class "glyphicon glyphicon-plus" ] []
            , text " Add new task"
            ]
        ]

      Just _ ->
        [ div
            [ class "btn-group" ]
            [ button
                [ class "btn btn-primary"
                , onClick address UberlistTaskList.SubmitNewTask
                ]
                [ span [ class "glyphicon glyphicon-plus" ] []
                , text " Submit new task"
                ]
            , button
                [ class "btn btn-warning"
                , onClick address UberlistTaskList.CancelNewTask
                ]
                [ span [ class "glyphicon glyphicon-remove" ] []
                , text " Add new task"
                ]
            ]
        ]
    )


main =
  StartApp.start
    { model =
        (UberlistTaskList.Model
          [ (UberlistTask.Model 1 "Hello world" False)
          , (UberlistTask.Model 2 "Hello world #2" False)
          ]
          0
          Nothing
        )
    , view = view
    , update = UberlistTaskList.update
    }
