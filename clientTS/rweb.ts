import WebSocket = require('ws');
class RWeb {
    address:string
    link:Map<string,Function>
    counter:number
    replyBind:Map<number,Function>
    conn:any
    constructor(address:string) {
        this.address = address
        this.link = new Map()
        this.counter = 0
        this.replyBind = new Map()
    }
    onmessage (e:any):void{
        let data = JSON.parse(e.data)
        if ('reply' in data && data.reply) {    //是对调用的回复
            let id = data.id;
            let f = this.replyBind.get(id)
            if (f == undefined) {
                console.log ("id:" + id + " is not exist")
                return
            }
            f (data.data)
            this.replyBind.delete(id)
        } else {    //是对本地的调用
            let f = this.link.get(data.function)
            if (f == undefined) {
                console.log("no function named:" + data.function + " linked")
                return
            }
            //开始构造调用
            let id = data.id
            let argument = data.argument
            let replier = new Replier(this,id)
            f (replier,...argument)
        }
    }
    async call(func: string,...args: any[]) :Promise<any>{
        // 开始发送调用
        let that = this;
        return new Promise(resolve => {
            let id = ++ that.counter
            console.log (id)
            that.replyBind.set(id,(data:any) => {
                console.log("Get@@@")
                resolve(data);
            })
            that.conn.send(JSON.stringify({
                'function':func,
                'argument':args,
                'id':id,
            }))
            setTimeout(()=>{
                this.replyBind.delete(id)
                //reject("timeout")
            },5000)
        })
    }
    onclose (e:any) {

    }
    onerror(e:any) {

    }
    async connect ():Promise<boolean> {
        this.conn = new WebSocket(this.address)
        //this.conn.that = this
        let that = this;
        this.conn.onmessage = (m:any) => {
            that.onmessage(m)
        }
        this.conn.onclose = this.onclose
        return new Promise(resolve => {
            this.conn.onopen = function () {
                resolve (true)
            }
        })
    }
}
class Replier {
    father:RWeb
    id:number
    call:Function
    replied:boolean
    constructor(father:RWeb,id:number) {
        this.father = father
        this.id = id
        this.call = father.call
        this.replied = false
    }
    reply(data:any):void{
        if (this.replied) {
            return
        }
        this.replied = true
        this.father.conn.send(JSON.stringify({
            "id":this.id,
            "reply":true,
            "data":data
        }))
    }
}
export {RWeb,Replier}