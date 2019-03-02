module Awoo exposing (Flags, Model, main)

import Dialog
import Html exposing (..)
import Html.Attributes as HA
import Html.Events as HE exposing (onClick, onInput)
import Json.Decode as JD
import Json.Decode.Pipeline exposing (decode, hardcoded, optional, required)
import Json.Encode as JE
import Time exposing (Time)
import WebSocket as WS


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
    , socketUrl : String
    , gameState : GameState
    , phase : Int
    , debugMsg : String
    , isLeader : Bool
    , role : Role
    , playerlist : List Player
    , rulesetName : String
    , roleDialog : Bool
    , voteDialog : Bool
    , targetedDialog : Bool
    , serverDialog : Bool
    , lastTargeted : TargetedDecoderContainer
    , votedUUID : String
    , serverMessage : String
    , tally : List TallyItem
    }


type alias Player =
    { name : String
    , uuid : String
    , leader : Bool
    }


type alias Role =
    { name : String
    , description : String
    , alive : Bool
    , team : Int
    , night_action : Int
    }


type GameState
    = Indeterminate
    | GameCannotBeJoined
    | NameEntry
    | Lobby
    | Day
    | Night
    | FinalGood
    | FinalEvil
    | Dead


emitGameState : GameState -> String
emitGameState state =
    case state of
        Day ->
            "day"

        Night ->
            "night"

        _ ->
            "how did i get here i am not good at computer"


type alias Flags =
    { socketUrl : String
    }


init : Flags -> ( Model, Cmd Msg )
init flags =
    ( { name = ""
      , socketUrl = flags.socketUrl
      , gameState = Indeterminate
      , phase = 0
      , debugMsg = ""
      , isLeader = False
      , role = Role "" "" False 0 0
      , playerlist = []
      , rulesetName = "Werewolf (X Players)"
      , roleDialog = False
      , voteDialog = False
      , targetedDialog = False
      , serverDialog = False
      , lastTargeted = TargetedDecoderContainer "" "" False
      , votedUUID = ""
      , serverMessage = ""
      , tally = []
      }
    , Cmd.none
    )



-- update


type Msg
    = Recv String
    | UpdateName String
    | EnterName
    | OpenRoleDialog
    | AckNonServerDialogs
    | AckServerDialog
    | OpenVoteDialog
    | VoteAction String


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Recv wsMsg ->
            case JD.decodeString messageDecoder wsMsg of
                Ok m ->
                    case m.message of
                        "awoo" ->
                            { model | gameState = NameEntry } ! []

                        "cannotjoin" ->
                            { model | gameState = GameCannotBeJoined } ! []

                        "pleasewait" ->
                            case JD.decodeValue pleasewaitDecoder m.payload of
                                Ok p ->
                                    { model | isLeader = p.isLeader, gameState = Lobby, debugMsg = wsMsg } ! [ WS.send model.socketUrl (JE.encode 0 (JE.object [ ( "pollPlayerList", JE.bool True ) ])) ]

                                Err err ->
                                    let
                                        _ =
                                            Debug.log "pleasewait payload error: " ( err, m.payload )
                                    in
                                    model ! []

                        "playerlist" ->
                            case JD.decodeValue playerlistDecoder m.payload of
                                Ok playerlist ->
                                    { model | playerlist = playerlist } ! []

                                Err err ->
                                    let
                                        _ =
                                            Debug.log "playerlist payload error" ( err, m.payload )
                                    in
                                    model ! []

                        "role" ->
                            case JD.decodeValue roleDecoder m.payload of
                                Ok role ->
                                    { model | role = role } ! []

                                Err err ->
                                    let
                                        _ =
                                            Debug.log "role payload error" ( err, m.payload )
                                    in
                                    model ! []

                        "phase" ->
                            case JD.decodeValue JD.int m.payload of
                                Ok phase ->
                                    let
                                        state =
                                            if phase % 2 == 1 then
                                                Day

                                            else
                                                Night
                                    in
                                    if model.gameState /= Dead then
                                        { model
                                            | gameState = state
                                            , phase = phase
                                            , voteDialog = False
                                            , roleDialog = False
                                            , targetedDialog = False
                                        }
                                            ! [ WS.send model.socketUrl (JE.encode 0 (JE.object [ ( "pollTally", JE.bool True ) ])) ]

                                    else
                                        model ! []

                                Err err ->
                                    let
                                        _ =
                                            Debug.log "gamestate payload error" ( err, m.payload )
                                    in
                                    model ! []

                        "tally" ->
                            case JD.decodeValue tallyDecoder m.payload of
                                Ok tally ->
                                    { model | tally = tally } ! [ WS.send model.socketUrl (JE.encode 0 (JE.object [ ( "pollPlayerList", JE.bool True ) ])) ]

                                Err err ->
                                    let
                                        _ =
                                            Debug.log "tally payload error" ( err, m.payload )
                                    in
                                    model ! []

                        "targeted" ->
                            case JD.decodeValue targetedDecoder m.payload of
                                Ok targetPayload ->
                                    -- inform the player concerning who died just now
                                    { model | targetedDialog = True, lastTargeted = targetPayload } ! []

                                Err err ->
                                    let
                                        _ =
                                            Debug.log "targeted payload error" ( err, m.payload )
                                    in
                                    model ! []

                        "dead" ->
                            { model | gameState = Dead } ! []

                        "privatemessage" ->
                            case JD.decodeValue JD.string m.payload of
                                Ok message ->
                                    { model | serverDialog = True, serverMessage = message } ! []

                                Err err ->
                                    let
                                        _ =
                                            Debug.log "private message payload error" ( err, m.payload )
                                    in
                                    model ! []

                        "victory" ->
                            case JD.decodeValue JD.int m.payload of
                                Ok victorypayload ->
                                    let
                                        victorystate =
                                            if victorypayload == 1 then
                                                FinalEvil

                                            else
                                                FinalGood
                                    in
                                    { model | gameState = victorystate } ! []

                                Err err ->
                                    let
                                        _ =
                                            Debug.log "victory payload error" ( err, m.payload )
                                    in
                                    model ! []

                        _ ->
                            let
                                _ =
                                    Debug.log "unknown command from server" wsMsg
                            in
                            { model | debugMsg = wsMsg } ! []

                Err err ->
                    let
                        _ =
                            Debug.log "recv payload error: " ( err, wsMsg )
                    in
                    model ! []

        UpdateName name ->
            { model | name = name } ! []

        EnterName ->
            model ! [ WS.send model.socketUrl (JE.encode 0 (JE.object [ ( "joinWithName", JE.string model.name ) ])) ]

        OpenRoleDialog ->
            { model | roleDialog = True } ! []

        AckNonServerDialogs ->
            { model | roleDialog = False, voteDialog = False, targetedDialog = False, votedUUID = "" } ! []

        AckServerDialog ->
            { model | serverDialog = False, serverMessage = "" } ! []

        OpenVoteDialog ->
            { model | voteDialog = True } ! []

        VoteAction player ->
            { model | voteDialog = False, votedUUID = player }
                ! [ WS.send model.socketUrl
                        (JE.encode 0
                            (JE.object
                                [ ( "voteFor", JE.string player )
                                , ( "time", JE.string (emitGameState model.gameState) )
                                ]
                            )
                        )
                  ]



-- decoders


type alias MessageDecoderContainer =
    { message : String
    , payload : JD.Value
    }


messageDecoder : JD.Decoder MessageDecoderContainer
messageDecoder =
    JD.map2 (\message payload -> MessageDecoderContainer message payload)
        (JD.field "message" JD.string)
        (JD.field "payload" JD.value)


type alias PleasewaitDecoderContainer =
    { isLeader : Bool }


pleasewaitDecoder : JD.Decoder PleasewaitDecoderContainer
pleasewaitDecoder =
    JD.map (\payload -> PleasewaitDecoderContainer payload)
        (JD.field "leader" JD.bool)


playerlistDecoder : JD.Decoder (List Player)
playerlistDecoder =
    JD.list playerDecoder


playerDecoder : JD.Decoder Player
playerDecoder =
    JD.map3 Player
        (JD.field "name" JD.string)
        (JD.field "uuid" JD.string)
        (JD.field "leader" JD.bool)


roleDecoder : JD.Decoder Role
roleDecoder =
    JD.map5 Role
        (JD.field "name" JD.string)
        (JD.field "description" JD.string)
        (JD.field "alive" JD.bool)
        (JD.field "team" JD.int)
        (JD.field "night_action" JD.int)


gamestateDecoder : JD.Decoder String
gamestateDecoder =
    JD.string


type alias TargetedDecoderContainer =
    { name : String
    , role : String
    , alive : Bool
    }


targetedDecoder : JD.Decoder TargetedDecoderContainer
targetedDecoder =
    JD.map3 TargetedDecoderContainer
        (JD.field "name" JD.string)
        (JD.field "role" JD.string)
        (JD.field "alive" JD.bool)


type alias TallyItem =
    { name : String
    , votes : List String
    }


tallyDecoder : JD.Decoder (List TallyItem)
tallyDecoder =
    JD.list
        (JD.map2 TallyItem
            (JD.at [ "candidate" ] JD.string)
            (JD.at [ "votes" ] (JD.list JD.string))
        )



-- subscriptions


subs : Model -> Sub Msg
subs model =
    WS.listen model.socketUrl Recv



--  view


view : Model -> Html Msg
view model =
    div [ HA.class "awoo-interface container" ]
        (case model.gameState of
            Indeterminate ->
                [ div [ HA.class "awoo-indeterminate row d-flex justify-content-center" ]
                    [ text "howling at the moon..." ]
                ]

            GameCannotBeJoined ->
                [ div [ HA.class "awoo-cannot-be-joined row d-flex justify-content-center" ]
                    [ text "game cannot presently be joined :(" ]
                ]

            NameEntry ->
                [ div [ HA.class "awoo-name-entry-dialog row d-flex justify-content-center" ]
                    [ text "OwO who's this?"
                    ]
                , div [ HA.class "awoo-name-entry-input row d-flex justify-content-center" ]
                    [ input [ HA.placeholder "player name", HA.value model.name, onInput UpdateName ] []
                    ]
                , div [ HA.class "awoo-name-entry-button row d-flex justify-content-center", onClick EnterName ]
                    [ i [ HA.class "fas fa-3x fa-check-circle" ] [] ]
                ]

            Lobby ->
                [ div [ HA.class "lobby row d-flex justify-content-center" ]
                    (renderLobby model)
                ]

            Day ->
                [ div [ HA.class "daytime row d-flex justify-content-center" ]
                    ([ renderGame model ] ++ [ dialog model ])
                ]

            Night ->
                [ div [ HA.class "nighttime row d-flex justify-content-center" ]
                    ([ renderGame model ] ++ [ dialog model ])
                ]

            Dead ->
                [ div [ HA.class "dead row d-flex justify-content-center" ]
                    [ div
                        [ HA.class "notification" ]
                        [ text "regrettably, you have died. whether you are a werewolf's snack, or fell prey to the machinations of paranoid villagers, surely someday you will have your revenge." ]
                    ]
                ]

            FinalGood ->
                [ div [ HA.class "finale-good row d-flex justify-content-center" ]
                    [ div [ HA.class "notification" ]
                        [ text "the game is over, and the good guys won! the villagers are safe... for now." ]
                    ]
                ]

            FinalEvil ->
                [ div [ HA.class "finale-evil row d-flex justify-content-center" ]
                    [ div [ HA.class "notification" ]
                        [ text "the game is over, and the werewolves won! the villagers are lunch... for now." ]
                    ]
                ]
        )


dialog : Model -> Html Msg
dialog model =
    Dialog.view
        (if model.roleDialog || model.voteDialog || model.targetedDialog then
            Just
                { closeMessage = Nothing
                , containerClass = Just "modal-container"
                , header =
                    Just
                        (div [ HA.class "d-flex fill", onClick AckNonServerDialogs ]
                            (if model.roleDialog then
                                [ roleDialogHeader, i [ HA.class "p-2 fa fa-times fa-4x cancel-icon" ] [] ]

                             else if model.voteDialog then
                                [ voteDialogHeader, i [ HA.class "p-2 fa fa-times fa-4x cancel-icon" ] [] ]

                             else if model.targetedDialog then
                                [ targetedDialogHeader, i [ HA.class "p-2 fa fa-times fa-4x cancel-icon" ] [] ]

                             else if model.serverDialog then
                                [ serverDialogHeader, i [ HA.class "p-2 fa fa-times fa-4x cancel-icon" ] [] ]

                             else
                                [ text "" ]
                            )
                        )
                , body =
                    Just
                        (div [ HA.class "modal-body" ]
                            (if model.roleDialog then
                                [ roleDialogBody model.role ]

                             else if model.voteDialog then
                                [ voteDialogBody model ]

                             else if model.targetedDialog then
                                [ targetedDialogBody model.lastTargeted ]

                             else if model.serverDialog then
                                [ serverDialogBody model.serverMessage ]

                             else
                                [ text "" ]
                            )
                        )
                , footer = Nothing
                }

         else if model.serverDialog then
            Just
                { closeMessage = Nothing
                , containerClass = Just "modal-container"
                , header =
                    Just
                        (div [ HA.class "d-flex fill", onClick AckServerDialog ]
                            [ serverDialogHeader, i [ HA.class "p-2 fa fa-times fa-4x cancel-icon" ] [] ]
                        )
                , body =
                    Just
                        (div [ HA.class "modal-body" ]
                            [ serverDialogBody model.serverMessage ]
                        )
                , footer = Nothing
                }

         else
            Nothing
        )


roleDialogHeader : Html Msg
roleDialogHeader =
    div [ HA.class "mr-auto p-2" ] [ span [ HA.class "modal-header-text" ] [ text "Your Role Is..." ] ]


voteDialogHeader : Html Msg
voteDialogHeader =
    div [ HA.class "mr-auto p-2" ] [ span [ HA.class "modal-header-text" ] [ text "Take Action On..." ] ]


targetedDialogHeader : Html Msg
targetedDialogHeader =
    div [ HA.class "mr-auto p-2" ] [ span [ HA.class "modal-header-text" ] [ text "Meanwhile..." ] ]


serverDialogHeader : Html Msg
serverDialogHeader =
    div [ HA.class "mr-auto p-2" ] [ span [ HA.class "modal-header-text" ] [ text "A sudden revelation!" ] ]


roleDialogBody : Role -> Html Msg
roleDialogBody role =
    div [ HA.class "d-flex fill flex-column" ]
        [ div [ HA.class "row role-name" ] [ text role.name ]
        , div [ HA.class "row role-description" ] [ text role.description ]
        , div [ HA.class "row role-image" ] [ img [ HA.src ("assets/" ++ role.name ++ ".jpg"), HA.class "img-fluid role-image" ] [] ]
        ]


targetedDialogBody : TargetedDecoderContainer -> Html Msg
targetedDialogBody lastTarget =
    div [ HA.class "d-flex fill flex-column" ]
        [ div [ HA.class "row last-targeted-name" ] [ text (lastTarget.name ++ " was marked for death!") ]
        , div []
            (if lastTarget.alive then
                [ div [ HA.class "row but-they-lived" ] [ text "...but they survived!" ] ]

             else
                [ div [ HA.class "row and-they-died" ] [ text ("They died, and in death the truth was revealed: They were a " ++ lastTarget.role ++ "!") ] ]
            )
        ]


voteDialogBody : Model -> Html Msg
voteDialogBody model =
    div [ HA.class "d-flex fill flex-column" ]
        [ ul [ HA.class "vote-list-select" ]
            (List.map renderVoteOption
                model.playerlist
            )
        ]


renderVoteOption : Player -> Html Msg
renderVoteOption player =
    li [ HA.class "vote-list-option", onClick (VoteAction player.uuid) ] [ text player.name ]


serverDialogBody : String -> Html Msg
serverDialogBody message =
    div [ HA.class "d-flex fill flex-column" ]
        [ div [ HA.class "row server-msg" ] [ text message ] ]


renderTally : List TallyItem -> List (Html Msg)
renderTally tally =
    [ ul [ HA.class "tally" ] (List.map renderTallyItem (List.sortBy votesInTallyItem tally) |> List.reverse) ]


votesInTallyItem : TallyItem -> Int
votesInTallyItem tally =
    List.length tally.votes


renderTallyItem : TallyItem -> Html Msg
renderTallyItem tally =
    let
        tallyCountString =
            List.length tally.votes |> toString
    in
    div [ HA.class "tallyitem" ]
        [ text (tally.name ++ "  (" ++ tallyCountString ++ ")")
        , div [ HA.class "tally-votes" ] (List.map renderTallyVote tally.votes)
        ]


renderTallyVote : String -> Html Msg
renderTallyVote vote =
    span [ HA.class "tally-vote" ] [ text vote ]


renderGame : Model -> Html Msg
renderGame model =
    div []
        [ div [ HA.class "interactions" ]
            [ button
                [ HA.type_ "button", HA.class "btn btn-info show-role", onClick OpenRoleDialog ]
                [ text "Show Role" ]
            , button
                [ HA.type_ "button", HA.class "btn btn-info show-vote", onClick OpenVoteDialog ]
                [ text "Take Action" ]
            ]
        , div [ HA.class "tally-list" ]
            (renderTally model.tally)
        ]


renderLobby : Model -> List (Html Msg)
renderLobby model =
    [ div []
        [ div
            [ HA.class "current-rules" ]
            [ span [] [ text "Current Ruleset:  " ]
            , if model.isLeader then
                span [] [ text model.rulesetName ]

              else
                span [] [ text model.rulesetName ]
            ]
        , div [ HA.class "playerlist" ]
            (renderLobbyList
                model.playerlist
            )
        ]
    ]


renderLobbyList : List Player -> List (Html Msg)
renderLobbyList players =
    [ ul [ HA.class "players" ] (List.map renderLobbyPlayer players) ]


renderLobbyPlayer : Player -> Html Msg
renderLobbyPlayer player =
    div [ HA.id ("player-" ++ player.uuid) ]
        [ text player.name
        ]
