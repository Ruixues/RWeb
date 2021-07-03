import {RWeb} from './rweb.js'

async function main() {
    let rweb = new RWeb("ws://127.0.0.1:1111/t")
    rweb.link.set("test", function (replier, t) {
        replier.reply("Ok")
    })
    await rweb.connect()
    rweb.call("test")
}

main()
