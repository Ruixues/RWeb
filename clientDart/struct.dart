class StandardCall {
  late String Function;
  late List<dynamic> Argument;
  late int Id;
  late bool IsReply;
  StandardCall(
      {required this.Function,
      required this.Argument,
      required this.Id,
      required this.IsReply}) {}
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
  late int Id;
  late bool Reply;
  late dynamic Data;
  StandardReply({required this.Id, required this.Reply, required this.Data}) {}
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
