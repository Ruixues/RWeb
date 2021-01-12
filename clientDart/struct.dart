class StandardCall {
  String Function;
  List<dynamic> Argument;
  int Id;
  bool IsReply;
  StandardCall({this.Function, this.Argument, this.Id, this.IsReply}) {}
  StandardCall.fromJson(Map<String, dynamic> json) {
    Function = json["function"];
    Argument = json["argument"];
    Id = json["id"];
    IsReply = json["reply"];
  }
  Map<String, dynamic> toJson() =>
      {'function': Function, 'argument': Argument, "id": Id, "reply": IsReply};
}

class StandardReply {
  int Id;
  bool Reply;
  dynamic Data;
  StandardReply({this.Id, this.Reply, this.Data}) {}
  Map<String, dynamic> toJson() {
    return {"id": Id, "reply": Reply, "data": Data};
  }

  StandardReply.fromJson(Map<String, dynamic> json) {
    Id = json["id"];
    Data = json["data"];
    Reply = json["reply"];
  }
}
/*
type StandardReply struct {
	Id    jsoniter.Number `json:"id"`
	Data  interface{}     `json:"data"`
	Reply bool            `json:"reply"`
}
*/
