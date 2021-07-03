import {Replier, RWeb} from "./rweb";

async function main() {
    let rweb = new RWeb("ws://127.0.0.1:1111/t")
    rweb.link.set("test", function (replier: Replier, t: string) {
        replier.reply("Hello Server!!!!!!!,You told me:" + t)
    })
    await rweb.connect()
    rweb.call("test")
}

main()