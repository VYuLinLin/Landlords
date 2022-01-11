import playerdata from "./data/player.js"
import eventlister from "./util/event_lister.js"

const myglobal = myglobal || {}
myglobal.playerData = playerdata()
window.$socket = eventlister({})

export default myglobal
