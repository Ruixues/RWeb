import 'dart:async';
import 'dart:convert';
import 'package:web_socket_channel/io.dart';

class _Call {
  final String function;
  final int id;
  final List<dynamic> arguments;
  _Call(this.function, this.arguments, this.id);
  Map toJson() {
    Map map = Map();
    map['function'] = function;
    map['argument'] = arguments;
    map['id'] = id;
    return map;
  }
}

class _Reply {
  final int id;
  final dynamic data;
  _Reply(this.id, this.data);
  Map toJson() {
    Map map = Map();
    map['id'] = id;
    map['data'] = data;
    map['reply'] = true;
    return map;
  }
}

typedef Reply(dynamic data);
typedef RBindFunction(Reply reply, List<dynamic> arguments);

class RWebClient {
  final String serverAddress;
  late IOWebSocketChannel conn;
  Map<String, RBindFunction> bindedFunction = Map();
  RWebClient({required this.serverAddress}) {}
  int idCount = 0;
  Map<int, Completer<dynamic>> link = Map();

  void onData(dynamic data) {
    if (!(data is String)) {
      return;
    }
    data = json.decode(data);
    // 开始解析
    if (data["reply"] == true) {
      // 是回复
      int? id = data["id"];
      if (!link.containsKey(id)) return;
      if (!link[id]!.isCompleted) {
        link[id]!.complete(data['data']);
      }
      return;
    }
    // 那就是函数调用
    String? functionName = data['function'];
    print("Call:" + (functionName ?? ""));
    if (functionName == null) return;
    List<dynamic> arguments = data['argument'] ?? List.empty();
    int? id = data['id'];
    if (id == null) return;
    if (bindedFunction.containsKey(functionName)) {
      //函数存在
      var f = bindedFunction[functionName]!;
      f((dynamic data) {
        _Reply reply = _Reply(id, data);
        conn.sink.add(json.encode(reply));
      }, arguments);
    }
  }

  bool bind(String function, RBindFunction f) {
    if (!bindedFunction.containsKey(function)) {
      bindedFunction[function] = f;
      return true;
    }
    return false;
  }

  void onError(dynamic error) {
    print(error);
    Future.delayed(Duration(seconds: 3), connect);
  }

  Future<bool> connect() async {
    conn = IOWebSocketChannel.connect(this.serverAddress);
    conn.stream.listen(onData, onError: onError);
    return true;
  }

  Future<dynamic> call(String function, List<dynamic> arguments) async {
    // 序列化
    var realId = ++idCount;
    var payload = _Call(function, arguments, realId);
    var data = jsonEncode(payload);
    conn.sink.add(data);
    var c = Completer<dynamic>();
    this.link[realId] = c;
    return c.future;
  }
}
