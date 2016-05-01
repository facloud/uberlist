module UberlistTask (..) where

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, on, targetValue, keyCode)
import Json.Decode as Json


type alias Model =
  { id : Int
  , title : String
  , titleIsField : Bool
  }


newTask =
  Model -1 "" True


type Action
  = EditTitle
  | UpdateTitleField String
  | SubmitNewTitle


update : Action -> Model -> Model
update action model =
  case action of
    EditTitle ->
      { model | titleIsField = True }

    UpdateTitleField t ->
      { model | title = t }

    SubmitNewTitle ->
      { model | titleIsField = False }


view : Signal.Address Action -> Model -> Html
view address model =
  div
    [ id "TaskDetails", class "col-lg-7" ]
    [ div
        [ class "page-header" ]
        [ viewTitle address model ]
    ]


onEnter : Signal.Address a -> a -> Html.Attribute
onEnter address value =
  on
    "keydown"
    (Json.customDecoder keyCode is13)
    (\_ -> Signal.message address value)


is13 : Int -> Result String ()
is13 code =
  if code == 13 then
    Ok ()
  else
    Err "not the right key code"


viewTitle : Signal.Address Action -> Model -> Html
viewTitle address model =
  if model.titleIsField then
    input
      [ class "form-control"
      , value model.title
      , placeholder "Create a new task"
      , on "input" targetValue (\title -> Signal.message address (UpdateTitleField title))
      , onEnter address SubmitNewTitle
      ]
      []
  else
    h3 [ onClick address EditTitle ] [ text model.title ]
