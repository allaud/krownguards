# README #

Standalone Legion TD game server (GO version)

### Install dependencies ###

    go get

### Run server ###

    go run main.go --usernames=Playername1,Playername2

### Generate differ ###

    cat models/*.go | python gen_diff.py > models/differ.go

### Test commands ###

    Add resources to player in slot (Slot required, other default = 0)
    http://localhost:8888/test?slot=1&gold=100&stone=100&food=10

    Change wave (Wave required)
    http://localhost:8888/test?wave=5

    Change king's health by value (King required: 0 = West, 1 = East, hp default = 0)
    http://localhost:8888/test?king=0&hp=1500

### Test server ###

    var serversocket = new WebSocket("ws://localhost:8888/ws");
    serversocket.onmessage = function(e) {console.log(e);};

    serversocket.send('{"type": "connect!", "username": "anton", "slot": 1}')

    serversocket.send('{"type": "pick_guild!", "slot": 1, "guild": "Warfactory", "uid": "ankxel"}')

    serversocket.send('{"type": "reset_guild!", "slot": 1, "uid": "ankxel"}')

    serversocket.send('{"type": "build!", "slot": 1, "coords": [5, 30], "name": "Goblin Technic", "uid": "ankxel"}')

    serversocket.send('{"type": "upgrade_unit!", "slot": 1, "id": 0, "upgrade": "Goblin Welder", "uid": "ankxel"}')

    serversocket.send('{"type": "sell_unit!", "slot": 1, "id": 0, "uid": "ankxel"}')

    serversocket.send('{"type": "upgrade_farm!", "slot": 1, "uid": "ankxel"}')
    serversocket.send('{"type": "cancel_farm_upgrade!", "slot": 1, "uid": "ankxel"}')

    serversocket.send('{"type": "upgrade_stone!", "slot": 1, "uid": "ankxel"}')
    serversocket.send('{"type": "cancel_stone_upgrade!", "slot": 1, "uid": "ankxel"}')

    serversocket.send('{"type": "send_income!", "slot": 1, "name": "Militia", "uid": "ankxel"}')

    serversocket.send('{"type": "upgrade_king!", "slot": 1, "attribute": "atk", "uid": "ankxel"}')
    serversocket.send('{"type": "upgrade_king!", "slot": 1, "attribute": "maxhp", "uid": "ankxel"}')
    serversocket.send('{"type": "upgrade_king!", "slot": 1, "attribute": "hpreg", "uid": "ankxel"}')

### Data send on connect ###

    Waves: {
        (int) Wave number: {
            (string) Unit:   "Unitname",
            (int)    Count:  15,
            (int)    Bounty: 13,
        }
    }

    TypeTooltips : {
        (string) Type : (string) Description,
    }

    AbilitiesTooltips : {
        (string) AbilityName : {
            (string) Type,
            (string) Range,
            (string) Cost,
            (string) Description,
        },
    }

    AffectedTooltips : {
        (string) EffectName : {
            (string) Type,
            (string) Description,
        },
    }

    GuildTooltips : {
        (string) GuildName : (string) Description,
    }