class StandardCall {
  public function: String = "";
  public argument: Array<any> = [];
  public id: String = "";
}
class rweb {
  private address: String = "";
  private conn: WebSocket | null = null;
  public onClose: Function | null = null;
  private bindFunction: Map<String, Function> = new Map();
  private bindReply: Map<String, Function> = new Map();
  private id: number = 0;
  constructor(address: String) {
    this.address = address;
  }
  private async onclose(): Promise<void> {
    // 准备重连
    this.conn = null;
    if (this.onClose != null) {
      this.onClose();
    }
    await this.sleep(2000); // 两秒后重连
    this.connect();
  }
  async sleep(ms: number) {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve("");
      }, ms);
    });
  }
  private onMessage(message: any): void {
    message = String(message);
    let json: any = JSON.parse(message);
    if (json.isReply) { //是回复
      let id: String = String(json.id);
      let data: any = json.data;
      if (this.bindReply.has(id)) {
        let f = this.bindReply.get(id);
        if (f != null) {
          f.call(data);
        }

        this.bindReply.delete(id);
      }
      return;
    }
  }
  public async call(name: String, ...args: any): Promise<any> {
    let c: StandardCall = new StandardCall();
    c.argument = args;
    c.function = name;
    c.id = (++this.id).toString();
    this.id %= 10000000000;
    let promise = new Promise((resolve) => {
      this.bindReply.set(c.id, resolve);
    });
    while (this.conn == null) {
        await this.sleep (1000);
    }
    this.conn.send(JSON.stringify(c));
    return promise;
  }
  public async connect(): Promise<void> {
    this.conn = new WebSocket(this.address.toString());
    this.conn.onclose = this.onclose;
    this.conn.onmessage = this.onMessage;
  }
}
export { rweb, StandardCall };
