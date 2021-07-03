class RWeb {
    constructor(address) {
        this.address = address
        this.link = new Map()
        this.counter = 0
        this.replyBind = new Map()
    }

    onmessage(e) {
        let data = JSON.parse(e.data)
        if ('reply' in data && data.reply) {    //是对调用的回复
            let id = data.id;
            if (!this.that.replyBind.has(id)) {
                console.log("id:" + id + " is not exist")
                return
            }
            this.that.relyBind.get(id)(data.data)
            this.that.replyBind.delete(id)
        } else {    //是对本地的调用
            let f = this.that.link.get(data.function)
            if (f == undefined) {
                console.log("no function named:" + data.function + " linked")
                return
            }
            //开始构造调用
            let id = data.id
            let argument = data.argument
            let replier = new Replier(this.that, id)
            f(replier, ...argument)
        }
    }

    async call(func, ...args) {
        // 开始发送调用
        return new Promise((resolve, reject) => {
            let id = ++this.counter
            this.replyBind.set(id, (data) => {
                resolve(data);
            })
            this.conn.send(JSON.stringify({
                'function': func,
                'argument': args,
                'id': id,
            }))
            setTimeout(5000, () => {
                this.replyBind.delete(id)
                reject("timeout")
            })
        })
    }

    onclose(e) {

    }

    onerror(e) {

    }

    async connect() {
        this.conn = new WebSocket(this.address)
        this.conn.that = this
        this.conn.onmessage = this.onmessage
        this.conn.onclose = this.onclose
        return new Promise((resolve, reject) => {
            this.conn.onopen = function () {
                resolve(true)
            }
            this.conn.onerror = (error) => {
                reject(error)
            }
        })
    }
}

class Replier {
    constructor(father, id) {
        this.father = father
        this.id = id
        this.call = father.call
        this.replied = false
    }

    reply(data) {
        if (this.replied) {
            return
        }
        this.replied = true
        this.father.conn.send(JSON.stringify({
            "id": this.id,
            "reply": true,
            "data": data
        }))
    }
}

export {RWeb, Replier}
