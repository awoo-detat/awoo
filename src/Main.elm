module Awoo exposing (Flags, Model, main)

import Html exposing (..)
import Html.Attributes as HA
import Html.Events as HE exposing (onClick, onInput)
import Json.Decode as JD
import Json.Decode.Pipeline exposing (decode, hardcoded, optional, required)
import Json.Encode as JE
import Time exposing (Time)
import Websocket as WS


main : Program Flags Model Msg
main =
    Html.programWithFlags
        { init = init
        , view = view
        , subscriptions = subs
        , update = update
        }


type alias Model =
    { name : String
    , response : String
    , socketUrl : String
    , currentTime : Time
    }


type alias Flags =
    { socketUrl : String
    }


init : Flags -> ( Model, Cmd Msg )
init flags =
    ( { name = "lolwat"
      , response = ""
      , socketUrl = flags.socketUrl
      , currentTime = 0
      }
    , Cmd.none
    )



-- update


type Msg
    = Tick Time
    | Recv String
    | UpdateName String


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Tick time ->
            { model | currentTime = time } ! []

        UpdateName name ->
            { model | name = name } ! []

        Recv wsMsg ->
            model ! []



-- decoders
-- subscriptions


subs : Model -> Sub Msg
subs model =
    Sub.batch [ Time.every Time.second Tick, WS.listen model.socketUrl Recv ]



--  view


view : Model -> Html Msg
view model =
    div [ class "awoo-interface" ]
        [ text "awoo"
        ]
